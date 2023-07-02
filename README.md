# GithubEQ
Syncs /bug reports in game to github issues.

![Issue In Game](http://i.imgur.com/blOKa6n.png)
- GMs can close and manage them inside github:

![Issue In Github](http://i.imgur.com/oxxs9uN.png)
- And see their details like this:
 
![Issue Detail in Github](http://i.imgur.com/HgksxNZ.png)


## Getting Started

Download from the [releases page](https://github.com/xackery/githubeq/releases) the githubeq binary. You can run it from anywhere that has access to the internet and your EQ database.

Run githubeq.exe, the first run it will say:
```
a new githubeq.conf file was created. Please open this file and configure githubeq, then run it again.
```

Edit the file. Some entries will require some additional steps:
- Database is straight forward, place your SQL connection settings.

- Github PersonalAccessToken:
	- [Click here](https://github.com/settings/personal-access-tokens/new) or manually go to your github profile, and click settings, developer settings, Personal Access Tokens, fine grained tokens.
	- Token name: GithubEQ
	- Expiration: custom, do it for a few years or whatever time you prefer.
	- Description: GithubEQ Token
	- Resource Owner: Set to the organization or repo you want this token to work at
	- Repository Access: Set to the repo you want issues to be managed at.
	- Repository Permissions: Issues (read and write)
	- Click Generate Token
	- Copy the token, and paste it into the PersonalAccessToken field in the githubeq.conf
	- Optional: If you are joining an organization, the request will be found in the Org's settings, Personal Access Token, Pending requests
- Github Labels: Go to the Issues tab of your repo, and click the `Labels` button. Create new labels that match the name of the labels in the above attributes, e.g.: `character, npc, item, ingame`. You can change these label names to your own custom ones, just be sure the labels exist, or you'll get errors later during the runtime.

On next run, you may see this message: `INF Successfully added github_issue_id column to bug_reports table`. This is normal, and means the database was updated to support github issues.