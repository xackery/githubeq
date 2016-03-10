# GithubEQ - Turn Everquest in game #issue command into a Github Issues tracking system!


Requirements
---
* The ability to compile the latest EQEMU (and do edits)
* Basic C++/Sql knowledge
* A github account


Edit eqemu_config.xml to add github entry.
---
* Add the following entries to eqemu_config.xml, near the bottom before </server>

```
<github personalaccesstoken="" repouser="" reponame="" issuelabel="ingame" characterlabel="character" npclabel="npc" itemlabel="item" refreshrate="120"/>
```
* Log in to github, and click your profile on the top right to access the drop down, and go to settings.
* Look for the `Personal access tokens` tab on the left sidebar menu.
* Click the `Generate new token` button on the top right.
* Enter your password to github.
* Enter in a generic description, like "Used for eq #issue system", and select the scopes `repo, repo:status, repo_deployment, public_repo`, and `read:repo_hook`.
* Save changes. You'll be shown the personalaccesstoken as a hash. Place the hash into the attribute noted above.
* Now set reponame/repouser. This is based on what repo you plan for issues to go to. For example, the repo these instructions are on is at the location https://github.com/Xackery/githubeq. repouser="Xackery", and reportname="githubeq" would be the fields for my repo. Change yours to your usage.
* Now go to the Issues tab of your repo, and click the `Labels` button. Create new labels that match the name of the labels in the above attributes, e.g.: `character, npc, item, ingame`. You can change these label names to your own custom ones, just be sure the labels exist, or you'll get errors later during the runtime.
* Refresh rate by default is every 2 minutes (120 seconds), you can change this to any frequency you like. 

Add database schema for issues
---
* Inside your database, copy the following SQL query inside your SQL client and run to create a new table:

```
CREATE TABLE `issues` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID for issue in game',
  `github_issue_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT 'issue id # inside github',
  `is_in_progress` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT 'flagged 1 when an issue is assigned to a person',
  `is_fixed` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT 'This is flagged 1 when an issue is closed.', 
  `is_deleted` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT 'When a user hits delete, this is flagged 1 and no longer shown to them.',
  `my_name` varchar(64) NOT NULL DEFAULT '',
  `my_account_id` int(10) unsigned NOT NULL,
  `my_character_id` int(10) unsigned NOT NULL,
  `my_zone_id` int(10) unsigned NOT NULL,
  `my_x` float NOT NULL,
  `my_y` float NOT NULL,
  `my_z` float NOT NULL,
  `message` varchar(512) NOT NULL DEFAULT '',
  `create_date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `last_modified` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `tar_name` varchar(64) NOT NULL DEFAULT '',
  `tar_is_npc` tinyint(1) NOT NULL,
  `tar_is_client` tinyint(1) NOT NULL,
  `tar_account_id` int(10) NOT NULL,
  `tar_character_id` int(10) NOT NULL,
  `tar_npc_type_id` int(10) NOT NULL,
  `tar_npc_spawngroup_id` int(11) NOT NULL,
  `item_id` int(11) NOT NULL,
  `item_name` varchar(64) NOT NULL DEFAULT '',
  `client` varchar(64) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=latin1;

```

Edit your source code to work with #issue
---
zone/command.cpp
Add the issue header
```
//Search for:
command_add("iplookup", "[charname] - Look up IP address of charname", 200, command_iplookup) ||

//Add this line after:
command_add("issue", "- Report an issue with the server", 0, command_issue) ||
```

zone/command.cpp
Add the command_issue function
```
//Search for:
void command_size(Client *c, const Seperator *sep)

//Add this code right BEFORE:
void command_issue(Client *c, const Seperator *sep) {

	if (sep->arg[1] && strcasecmp(sep->arg[1], "delete") == 0) { //Delete an issue
		if (!sep->arg[2] || atoi(sep->arg[2]) == 0) {
			c->Message(0, "Invalid issue id. Format: #issue delete <number>");
			return;
		}
		uint32 issue_id = atoi(sep->arg[2]);
		std::string query = StringFormat("UPDATE issues SET is_deleted = 1 WHERE is_deleted = 0 AND my_character_id = %u AND id = %u LIMIT 1", c->CharacterID(), issue_id);
		auto results = database.QueryDatabase(query);
		if (!results.Success()) {
			c->Message(13, "Deleting an issue failed. The admins have been notified.");
			Log.Out(Logs::General, Logs::Normal, "#issue creation failed for user %u: %s", c->CharacterID(), results.ErrorMessage().c_str());
			return;
		}
		c->Message(0, "You have deleted issue #%u.", issue_id);
		return;
	}

	if (sep->arg[1] && strcasecmp(sep->arg[1], "list") == 0) { //List past issues
		//List Issues
		std::string query = StringFormat("SELECT id, github_issue_id, is_in_progress, is_fixed, tar_name, message FROM issues WHERE my_character_id = %u AND is_deleted = 0 ORDER BY last_modified DESC LIMIT 10", c->CharacterID());
		auto results = database.QueryDatabase(query);
		if (!results.Success()) {
			c->Message(13, "Listing issues failed. The admins have been notified.");
			Log.Out(Logs::General, Logs::Normal, "#issue list failed for user %u: %s", c->CharacterID(), results.ErrorMessage().c_str());
			return;
		}
		if (results.RowCount() == 0) {
			c->Message(0, "You have no pending issues.");
			return;
		}

		c->Message(0, "Your %u most recently updated issues:", results.RowCount());
		for (auto row = results.begin(); row != results.end(); ++row) {
			std::string status = "New";
			if (atoi(row[1]) > 0) status = "Reported";
			if (atoi(row[2]) == 1) status = "In Progress";
			if (atoi(row[3]) == 1) status = "Fixed";

			std::string details = "";
			if (strlen(row[4]) > 0 && strcasecmp(row[4], "(null)") != 0) details.append(StringFormat("(%s) ", row[4]));
			std::string deletecommand = StringFormat("#issue delete %u", atoi(row[0]));
			details.append(StringFormat("%s", row[5]));

			c->Message(0, "#%u status: %s, details: %s [ %s ]", atoi(row[0]), status.c_str(), details.c_str(), c->CreateSayLink(deletecommand.c_str(), "delete").c_str());
		}
		return;
	}


	if (!sep->arg[1] || (strlen(sep->arg[1]) == 0)) {
		uint32 issue_count = 0;
		std::string query = StringFormat("SELECT count(id) FROM issues WHERE my_character_id = %u AND is_deleted = 0 LIMIT 1", c->CharacterID());
		auto results = database.QueryDatabase(query);
		if (results.Success()) {
			if (results.RowCount() == 1) {
				auto row = results.begin();
				issue_count = atoi(row[0]);
			}
		}		

		if (issue_count > 0) c->Message(0, "You have %u previously submitted issues. [ %s ]", issue_count, c->CreateSayLink("#issue list", "list").c_str());
		c->Message(0, "To report something to the GMs, you may target a mob or player and then:");
		c->Message(0, "/say #issue Your report message");
		return;
	}
	

	std::string itemname = "";
	uint32 itemid = 0;
	auto inst = c->GetInv()[MainCursor];
	if (inst) { 
		auto item = inst->GetItem();
		if (item) {
			itemname = item->Name;
			itemid = item->ID;
		}
	}

	std::string clientversion = "";
	
	if (c->GetClientVersion() == ClientVersion::Titanium) clientversion = "Titanium";
	else if (c->GetClientVersion() == ClientVersion::SoF) clientversion = "SoF";
	else if (c->GetClientVersion() == ClientVersion::SoD) clientversion = "SoD";
	else if (c->GetClientVersion() == ClientVersion::UF) clientversion = "Underfoot";
	else if (c->GetClientVersion() == ClientVersion::UF) clientversion = "UF";
	else if (c->GetClientVersion() == ClientVersion::RoF) clientversion = "RoF";
	else if (c->GetClientVersion() == ClientVersion::RoF2) clientversion = "RoF2";
	else clientversion = "Unknown";
	
	std::string query = StringFormat("INSERT INTO issues"
		"(my_name, my_account_id, my_character_id, my_zone_id, my_x, my_y, my_z, message, tar_name, tar_is_npc, tar_is_client, tar_account_id, tar_character_id, tar_npc_type_id, tar_npc_spawngroup_id, item_id, item_name, client)"
		"VALUES (\"%s\", %u, %u, %u, %f, %f, %f, \"%s\", \"%s\", %u, %u, %u, %u, %u, %u, %u, \"%s\", \"%s\")",
		EscapeString(c->GetName()).c_str(),
		c->AccountID(),
		c->CharacterID(),
		c->GetZoneID(),
		c->GetX(),
		c->GetY(),
		c->GetZ(),
		EscapeString(sep->argplus[1]).c_str(), //message
		((c->GetTarget() == nullptr) ? 0 : EscapeString(c->GetTarget()->GetCleanName()).c_str()),
		((c->GetTarget() == nullptr || !c->GetTarget()->IsNPC()) ? 0 : 1),
		((c->GetTarget() == nullptr || !c->GetTarget()->IsClient()) ? 0 : 1),
		((c->GetTarget() == nullptr || !c->GetTarget()->IsClient()) ? 0 : c->GetTarget()->CastToClient()->AccountID()),
		((c->GetTarget() == nullptr || !c->GetTarget()->IsClient()) ? 0 : c->GetTarget()->CastToClient()->CharacterID()),
		((c->GetTarget() == nullptr || !c->GetTarget()->IsNPC()) ? 0 : c->GetTarget()->CastToNPC()->GetNPCTypeID()),		
		((c->GetTarget() == nullptr || !c->GetTarget()->IsNPC()) ? 0 : c->GetTarget()->CastToNPC()->GetSp2()),
		itemid,
		EscapeString(itemname.c_str()).c_str(),
		EscapeString(clientversion.c_str()).c_str()
		);
	auto results = database.QueryDatabase(query);
	if (!results.Success()) {
		c->Message(13, "Creating an issue failed. The admins have been notified.");
		Log.Out(Logs::General, Logs::Normal, "#issue creation failed for user %u: %s", c->CharacterID(), results.ErrorMessage().c_str());
		return;		
	}

	c->Message(0, "Your issue #%i has been submitted.", results.LastInsertedID());
}
```

zone/command.h
```
Search for:
void command_iplookup(Client *c, const Seperator *sep);

Add After:
void command_issue(Client *c, const Seperator *sep);
```

* Compile your source code, and run. In game, you should now be able to type #issue, and report something with #issue SomeText Here.


Run githubeq.exe
---
In the same directory your eqemu_config.xml file resides, run githubeq.exe (or githubeq if linux/osx) and it should see the github config options.
