from discord.ext import commands
from discord import Embed, File


def setup(bot):
    bot.add_cog(ErrorHandler(bot))

class CannotConnectToAPI(Exception):
    pass

class ErrorHandler(commands.Cog):
    def __init__(self, bot):
        self.bot = bot

    @commands.Cog.listener()
    async def on_command_error(self, ctx, error):
        if hasattr(error, 'on_error'):
            return

        elif isinstance(error, commands.MissingPermissions):

            if not ctx.guild.me.guild_permissions.embed_links:
                return await ctx.send("You need to be an administrator to run this command")
            error_embed = Embed(color=0xFF8C00)
            file = File("static/fist.png", filename="fist.png")
            error_embed.set_thumbnail(url="attachment://fist.png")
            error_embed.title = "Missing Permissions"
            error_embed.description = "You need to be an administrator to run this command"
            return await ctx.send(embed=error_embed, file=file)

        elif isinstance(error, commands.BotMissingPermissions):
            await ctx.send(f'{self.bot.user.mention} needs permission to send Embeds and create Invite Links')
            if ctx.guild.me.guild_permissions.embed_links:
                file = File("static/perms.png", filename="perms.png")
                return await ctx.send(file=file)
        else:
            raise error
        """
        elif isinstance(error, commands.CommandInvokeError):
            api_embed = Embed(color=0x006994)
            file = File("static/fist.png", filename="fist.png")
            api_embed.set_thumbnail(url="attachment://fist.png")
            return await ctx.send(embed=api_embed, file=file)
            return
            """


