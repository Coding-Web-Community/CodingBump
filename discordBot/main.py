from discord.ext import commands
from discord import Embed, File, TextChannel, Status, CustomActivity, ActivityType, Invite
from datetime import timedelta
from discord.errors import LoginFailure
from dotenv import load_dotenv
from os import getenv, listdir
from aiohttp import ClientSession
from json import dumps, loads
from datetime import datetime, timedelta


class CannotConnectToAPI(Exception):
    pass


# Guild Cache
class CacheService(object):

    def __init__(self):
        self.data = {}

    def __setitem__(self, key: int, item: int):
        self.data[key] = item

    def __getitem__(self, key):
        return self.data[key]


GuildCache = CacheService()


# Custom Cooldown Mapping
class CustomBumpCooldown:
    def __init__(self, rate: int, per: float, bucket: commands.BucketType):
        self.default_mapping = commands.CooldownMapping.from_cooldown(rate, per, bucket)

    def __call__(self, ctx: commands.Context):
        bucket = self.default_mapping.get_bucket(ctx.message)
        retry_after = bucket.update_rate_limit()
        # check if TooEarly but not in cache
        print(GuildCache.data)
        if ctx.guild.id not in GuildCache.data:
            GuildCache[ctx.guild.id] = -1
        if GuildCache[ctx.guild.id] == -1:
            return True
        if retry_after:
            raise commands.CommandOnCooldown(bucket, retry_after)
        return True


# POST request to API (bump)
async def ready_for_bump(guildId: int):
    payload = {'guildId': guildId}
    try:
        async with ClientSession() as session:
            async with session.post("http://localhost:8080/V1/bump", data=dumps(payload)) as resp:
                data = loads(await resp.text())
                print(data['code'], data['payload'])
                return data['code'], data['payload']
    except Exception:
        raise CannotConnectToAPI


# Re-add guild to Cache if not exist but on cooldown, else replace timestamp
async def handle_status(guildId: int, status: int, data: dict):
    bumped = datetime.fromtimestamp(data['timestamp'])
    now = datetime.now()
    if status == 200:
        GuildCache.data[guildId] = int(datetime.timestamp(now))
        print(GuildCache.data)
        return status, None
    if GuildCache.data[guildId] == -1:
        to_bump = bumped + timedelta(seconds=60)
        remain = to_bump - now
        if status == 425:
            return status, remain.total_seconds()


load_dotenv()
# load the .env file

bot = commands.Bot(
    command_prefix=f"{getenv('BOT_PREFIX') or '`'}"
)

# for cogs
for cog in listdir('cog'):
    if '.py' not in cog:
        continue
    print(f"Loaded cog: {cog[:-3]}")
    bot.load_extension(f'cog.{cog[:-3]}')


@bot.listen('on_ready')
async def print_stats():
    print(f"Logged in as: {bot.user}")
    print(f"Shard-id: {bot.shard_id or 0}")
    print(f"Shards: {bot.shard_count or 0}")
    await bot.change_presence(status=Status.online,
                              activity=CustomActivity(type=ActivityType.listening,
                                                      name=f'to {bot.command_prefix}bump'))


# Bump command
@bot.command(brief="Bump the server",
             help="Use it to bump your server every 2hrs",
             usage='`bump',
             name='bump')
@commands.max_concurrency(number=1, per=commands.BucketType.guild, wait=False)
@commands.guild_only()
@commands.bot_has_permissions(embed_links=True)
@commands.check(CustomBumpCooldown(1, 60.0, commands.BucketType.guild))
async def _bump(ctx):
    # Bumping logic handled by cooldown mapping
    print("HERE")
    status, data = await ready_for_bump(guildId=ctx.guild.id)
    status, retry = await handle_status(guildId=ctx.guild.id, status=status, data=data)
    if status == 425:
        file, embed = await cooldown_error(retry_after=retry)
        embed.timestamp = ctx.message.created_at
        await ctx.send(embed=embed, file=file)
        ctx.reset_cooldown()
        return

    bump_embed = Embed(color=0xCCCC00)
    bump_embed.title = "Bumped"
    bump_embed.url = "https://blacksmithop.xyz/"  # points to guild in bot list
    bump_embed.set_thumbnail(url=bot.user.avatar_url)
    bump_embed.timestamp = ctx.message.created_at
    bump_embed.set_footer(text=f"Guild ID: {ctx.guild.id}", icon_url=ctx.guild.icon_url)
    bump_embed.description = f"{ctx.author.mention} bumped {ctx.guild.name}"
    await ctx.send(embed=bump_embed)
    ctx.reset_cooldown()


# Bump errors
@_bump.error
async def show_remaining(ctx, error):
    if isinstance(error, commands.CommandOnCooldown):
        file, embed = await cooldown_error(retry_after=error.retry_after)
        embed.timestamp = ctx.message.created_at
        await ctx.send(embed=embed, file=file)


# Create embed showing Bump cooldown
async def cooldown_error(retry_after: int):
    bump_embed = Embed(color=0xcf142b)
    bump_embed.title = "Bump on cooldown"
    bump_embed.url = "https://blacksmithop.xyz/"  # points to guild in bot list
    bump_embed.set_thumbnail(url=bot.user.avatar_url)

    time_left = str(timedelta(0, seconds=retry_after))[:-7]
    bump_embed.description = f"You can bump again in **{time_left}**"
    file = File("static/stop.png", filename="stop.png")
    bump_embed.set_thumbnail(url="attachment://stop.png")
    return file, bump_embed


bot.remove_command('help')


# Simple help embed
@bot.command(name='help')
@commands.bot_has_permissions(embed_links=True)
async def _help(ctx):
    help_embed = Embed()
    file = File("static/think.png", filename="think.png")
    help_embed.set_thumbnail(url="attachment://think.png")
    help_embed.title = "**Commands**"
    help_embed.description = "View the [Bot List](https://blacksmithop.xyz/)"
    help_embed.add_field(name="\u200b",
                         value="```fix\nhelp``` ```diff\n- show this message```", inline=False)
    help_embed.add_field(name="\u200b",
                         value="```fix\nbump``` ```diff\n- bump this server```", inline=False)
    help_embed.add_field(name="\u200b",
                         value="```fix\ninvite``` ```diff\n- set or view the current invite channel```")
    await ctx.send(embed=help_embed, file=file)


# Set invite channel, fetch invite by bot from channel if set to never expires
@bot.command(name='invite')
@commands.bot_has_permissions(create_instant_invite=True, embed_links=True)
@commands.has_permissions(administrator=True)
async def _invite(ctx, channel: TextChannel = None):
    # check if channel invite is in API when no channel is passed
    if channel is None:
        channel = ctx.channel
    # User wishes to set current / argument channel as the invite channel
    # Check is an invite made by bot exists for channel
    invites = await channel.invites()
    invite: Invite
    usable_invite = next((invite for invite in invites
                          if not (invite.max_age and invite.max_uses) and invite.inviter == bot.user),
                         None)
    if usable_invite is None:
        usable_invite = await channel.create_invite(max_age=0, max_uses=0,
                                                    temporary=False)
    invite_embed = Embed()
    file = File("static/raise.png", filename="raise.png")
    invite_embed.set_thumbnail(url="attachment://raise.png")
    invite_embed.title = "Invite Channel"
    invite_embed.url = "https://blacksmithop.xyz/"
    invite_embed.description = f"Current invite channel is set to {usable_invite.channel.mention}" \
                               f"\n[Invite Url]({usable_invite.url})"

    await ctx.send(embed=invite_embed, file=file)


# check ping
@bot.command(name='ping')
@commands.bot_has_permissions(embed_links=True)
async def _ping(ctx):
    ping_embed = Embed()
    file = File("static/ping.png", filename="raise.png")
    ping_embed.set_thumbnail(url="attachment://raise.png")
    ping_embed.add_field(name='Discord ', value=f'{round(bot.latency,3)}s')
    await ctx.send(embed=ping_embed, file=file)

if __name__ == '__main__':
    token = getenv('BOT_TOKEN')
    try:
        bot.run(token)
    except LoginFailure:
        print("Improper Token was passed!")
