## Github EQ

Creates a new command `#issue` in EverQuest that players can use as a "fast way to give feedback".

Usage: `#issue <msg that can be as long as desired>`

### Installation

1. Go to the [releases page](https://github.com/xackery/githubeq/releases) or click links in instructions below.
1. Download [githubeq-linux](https://github.com/xackery/githubeq/releases/latest/download/githubeq-linux) or [githubeq-windows.exe](https://github.com/xackery/githubeq/releases/latest/download/githubeq-windows.exe) based on your OS, place it in the root of your eqemu folder. (Ideally below quests)
1. Run githubeq. It will exit and generate a githubeq.conf file.
1. Create a personal access token:
    1. Go to [Settings in Github](https://github.com/settings/profile)
    1. On the bottom, click [Developer Settings](https://github.com/settings/apps)
    1. Go to [Fine-grained tokens](https://github.com/settings/tokens?type=beta)
    1. Click [Generate new Token](https://github.com/settings/personal-access-tokens/new)
    1. For token name, put `githubeq-<servername>` or whatever you like. 
    1. Expiration, I set mine for a year. This means in a year you'll need to refresh this token. (Not sure the maximum duration it can go)
    1. Repository Access, only select repositories, and find your repository you want issues to go to.
    1. Permissions, expand and scroll down to Issues, set access to Read and Write 
    1. Click generate token. You'll get a token with prefix like `github_pat_` and a long string of characters. Copy this token.
1. Edit your githubeq.conf file and put in your token.
1. Add repository and user based on what the token was used, e.g. user: xackery, repository: githubeq would be for this repository.
1. Run githubeq again. It should now work and go idle.
1. Check if githubeq made an `issues` folder. This is where the issues in lua should generate at.
1. Download [issue.lua](https://github.com/xackery/githubeq/releases/latest/download/issue.lua) and place it in your `quests/lua_modules/commands` folder.
1. Edit your `quests/lua_modules/command.lua` file and add the following line to the below the other 4 or 5 commands: `commands["issue"] 	  = { 0,   require(commands_path .. "issue") };`
1. Go in game and type `#rq` to reload quests.
1. Type `#issue Hello GithubEQ!` in game and you should get a message your issue is submitted.
1. By default, check_frequency_minutes is set to 1, meaning every minute it'll look for issues in the issues folder, and push them to your github Issues page. Wait 60 seconds and verify all is okay.

