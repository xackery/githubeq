---@param e ModRegisterBug
function RegisterBug(e)
    local msg = ""
    local client = e.self
    local inv = client:GetInventory()

    msg = msg .. "-------\n"
    msg = msg .. "message\n"
    msg = e.bug_report .. "\n"

    msg = msg .. "-------\n"
    msg = msg .. "preview\n"

    msg = msg .. "Cursor: "
    if inv:GetItem(Slot.Cursor).valid then
        msg = string.format("%s%s (%d)\n", msg, inv:GetItem(Slot.Cursor):GetName(), inv:GetItem(Slot.Cursor):GetID())
    else
        msg = string.format("%snone\n", msg)
    end

    msg = msg .. string.format("Location: `#zone %s %d %d %d`\n", e.zone,  client:GetX() , client:GetY() , client:GetZ())

    msg = msg .. "-------\n"
    msg = msg .. "bug info\n"

    msg = msg .. "zone: " .. e.zone .. "\n"
    msg = msg .. "client_version_id: " .. e.client_version_id .. "\n"
    msg = msg .. "client_version_name: " .. e.client_version_name .. "\n"
    msg = msg .. "account_id: " .. e.account_id .. "\n"
    msg = msg .. "character_id: " .. e.character_id .. "\n"
    msg = msg .. "character_name: " .. e.character_name .. "\n"
    msg = msg .. "reporter_spoof: " .. e.reporter_spoof .. "\n"
    msg = msg .. "category_id: " .. e.category_id .. "\n"
    msg = msg .. "category_name: " .. e.category_name .. "\n"
    msg = msg .. "reporter_name: " .. e.reporter_name .. "\n"
    msg = msg .. "ui_path: " .. e.ui_path .. "\n"
    msg = msg .. "pos_x: " .. e.pos_x .. "\n"
    msg = msg .. "pos_y: " .. e.pos_y .. "\n"
    msg = msg .. "pos_z: " .. e.pos_z .. "\n"
    msg = msg .. "heading: " .. e.heading .. "\n"
    msg = msg .. "time_played: " .. e.time_played .. "\n"
    msg = msg .. "target_id: " .. e.target_id .. "\n"
    msg = msg .. "target_name: " .. e.target_name .. "\n"
    msg = msg .. "optional_info_mask: " .. e.optional_info_mask .. "\n"
    msg = msg .. "_can_duplicate: " .. e._can_duplicate .. "\n"
    msg = msg .. "_crash_bug: " .. e._crash_bug .. "\n"
    msg = msg .. "_target_info: " .. e._target_info .. "\n"
    msg = msg .. "_character_flags: " .. e._character_flags .. "\n"
    msg = msg .. "_unknown_value: " .. e._unknown_value .. "\n"
    msg = msg .. "system_info: " .. e.system_info .. "\n"

    msg = msg .. "-------\n"
    msg = msg .. "character info\n"

    msg = msg .. string.format("Name: %s\n", client:GetName())
    msg = msg .. string.format("Account: %s\n", client:AccountName())
    msg = msg .. string.format("Time: %s\n", os.date())
    msg = msg .. string.format("Level: %d\n", client:GetLevel())
    msg = msg .. string.format("Class: %s\n", client:GetClass())

    msg = msg .. "-------\n"

    msg = msg .. "inventory\n"

    --- a dictionary of slots
    local slots = {
        [Slot.Primary] = "Primary",
        [Slot.Secondary] = "Secondary",
        [Slot.Charm] = "Charm",
        [Slot.Ear1] = "Ear1",
        [Slot.Head] = "Head",
        [Slot.Face] = "Face",
        [Slot.Ear2] = "Ear2",
        [Slot.Neck] = "Neck",
        [Slot.Shoulders] = "Shoulders",
        [Slot.Arms] = "Arms",
        [Slot.Back] = "Back",
        [Slot.Wrist1] = "Wrist1",
        [Slot.Wrist2] = "Wrist2",
        [Slot.Range] = "Range",
        [Slot.Hands] = "Hands",
        [Slot.Finger1] = "Finger1",
        [Slot.Finger2] = "Finger2",
        [Slot.Chest] = "Chest",
        [Slot.Legs] = "Legs",
        [Slot.Feet] = "Feet",
        [Slot.Waist] = "Waist",
        [Slot.PowerSource] = "PowerSource",
        [Slot.Ammo] = "Ammo",
        [Slot.Cursor] = "Cursor",
        [Slot.Shoulder] = "Shoulder",
        [Slot.Bracer1] = "Bracer1",
        [Slot.Bracer2] = "Bracer2",
        [Slot.Ring1] = "Ring1",
        [Slot.Ring2] = "Ring2",
    }

    for slot, name in pairs(slots) do
        local item = inv:GetItem(slot)
        if item ~= nil and item.valid then
            msg = string.format("%s%s: %s (%d)\n", msg, name, item:GetName(), item:GetID())
        else
            msg = string.format("%s%s: none\n", msg, name)
        end
    end

    for i = Slot.PossessionsBegin, Slot.PossessionsEnd do
        local item = inv:GetItem(i)
        if item ~= nil and item.valid then
            msg = string.format("%sPossession %d: %s (%d)\n", msg, i, item:GetName(), item:GetID())
        else
            --msg = string.format("%sPossession %d: none\n", msg, i)
        end
    end

    for i = Slot.PossessionsBagsBegin, Slot.PossessionsBagsEnd do
        local item = inv:GetItem(i)
        if item ~= nil and item.valid then
            msg = string.format("%sPossession Bag %d: %s (%d)\n", msg, i, item:GetName(), item:GetID())
        else
            --msg = string.format("%sPossession Bag %d: none\n", msg, i)
        end
    end


    for i = Slot.CursorBagBegin, Slot.CursorBagEnd do
        local item = inv:GetItem(i)
        if item ~= nil and item.valid then
            msg = string.format("%sCursor Bag %d: %s (%d)\n", msg, i, item:GetName(), item:GetID())
        else
            --msg = string.format("%sCursor Bag %d: none\n", msg, i)
        end
    end

    local item = inv:GetItem(Slot.Tradeskill)
    if item ~= nil and item.valid then
        msg = string.format("%sTradeskill: %s (%d)\n", msg, item:GetName(), item:GetID())
    else
        msg = string.format("%sTradeskill: none\n", msg)
    end

    for slot = Slot.BankBagsBegin, Slot.BankBagsEnd do
        local item = inv:GetItem(slot)
        if item ~= nil and item.valid then
            msg = string.format("%sBank Bag %d: %s (%d)\n", msg, slot, item:GetName(), item:GetID())
        else
            --msg = string.format("%sBank Bag %d: none\n", msg, slot)
        end
    end

    for slot = Slot.SharedBankBagsBegin, Slot.SharedBankBagsEnd do
        local item = inv:GetItem(slot)
        if item ~= nil and item.valid then
            msg = string.format("%sShared Bank Bag %d: %s (%d)\n", msg, slot, item:GetName(), item:GetID())
        else
            --msg = string.format("%sShared Bank Bag %d: none\n", msg, slot)
        end
    end

    local filename = "issues/" .. e.self:AccountName() .. "-" .. math.random(1000000, 9999999) .. ".txt";

    local w, err = io.open(filename, "a");
    if not w then
        e.self:Message(MT.White, "failed opening issue: " .. err);
        return
    end
    _, err = w:write(msg)
    if err then
        e.self:Message(MT.White, "failed writing issue: " .. err);
        return
    end
    --file:write(e.self:GetName() .. " - " .. e.self:GetAccountName() .. " - " .. e.self:GetIP() .. " - " .. os.date() .. " - " .. msg .. "\n");
    local isClosed, _, code = w:close();
    if not isClosed then
        e.self:Message(MT.White, "Error closing issue file: " .. code);
        return
    end

    e.IgnoreDefault = true
    return e
end
