## Github EQ

Uses a lua mod to redirect /bug reports to github issues.

### What does it look like?

When someone does a /bug in game, it sends the issue to github. No database interactions.

How the issues look:

<img src="https://github.com/xackery/githubeq/assets/845670/53d2d5f4-fa34-4e67-92ef-adf44378b323" width="200">

Expanded sections:

<img src="https://github.com/xackery/githubeq/assets/845670/19e685bb-8be3-42bd-872b-815e5b83e12a" width="200">

### How does it work?

- On /bug, a lua mod is fired and overrides the default behavior of in game #bugs.
- The lua mod creates a .txt file in the issues folder on your eqemu dir.
- Githubeq is a program that every minute by default, iteartes txt files in issues, and sends it to github as an issue. It then deletes the .txt file.
- You can then use projects in github to organize the issues and triage things, with less noise like you get on discord forum bug reports.

### Installation

1. *NOTE* eqemu has not yet added [a PR](https://github.com/EQEmu/Server/pull/4209) that adds support for this. Once it is added, you'll need to wait for a release from eqemu and update your binaries to obtain it.
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
1. Download [register_bug.lua](https://github.com/xackery/githubeq/releases/latest/download/register_bug.lua) and place it in your `mods` folder, which should be in your eqemu root.
1. Create a file called `mods/load_order.txt` if it doesn't already exist.
1. Edit `mods/load_order.txt` and add `register_bug.lua` to the bottom of the file.
1. Go in game and type `#rq` to reload quests.
1. Type `/bug` in game and create an issue.
1. Peek at githubeq's output after ~60s and you should see if it picked up the file.
1. Check if issues in your repository is populated with your new issue.