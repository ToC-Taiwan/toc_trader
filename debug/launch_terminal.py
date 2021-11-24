#!/usr/bin/python3
import iterm2
import AppKit
import pathlib

# Launch the app
AppKit.NSWorkspace.sharedWorkspace().launchApplication_("iTerm2")


async def main(connection):
    app = await iterm2.async_get_app(connection)
    # Foreground the app
    await app.async_activate()
    path = pathlib.Path(__file__).parent.absolute().__str__()+'/debugee.sh'
    # This will run 'vi' from bash. If you use a different shell, you'll need
    # to change it here. Running it through the shell sets up your $PATH so you
    # don't need to specify a full path to the command.
    # await iterm2.Window.async_create(connection, command="/bin/bash -l -c vi")
    await iterm2.Window.async_create(connection, command=path)


# Passing True for the second parameter means keep trying to
# connect until the app launches.
iterm2.run_until_complete(main, True)
