--- @param e PlayerEventCommand
local function issue(e)
    if #e.args < 2 then
        e.self:Message(MT.White, "usage: #issue <msg> - Items held, mobs targetted, are all included in the issue report");
        return
    end

    local client = e.self

    local msg = client:AccountName() .. "\n"
    tags = ""
    local inv = client:GetInventory()


    if inv:GetItem(Slot.Cursor).valid then
        tags = tags .. "item,"
    end

    if client:GetTarget().valid then
        tags = tags .. "target,"
    end

    if tags == "" then
        tags = "none"
    end

    msg = msg .. tags .. "\n"
    msg = msg .. table.concat(e.args, " ", 1) .. "\n";

    msg = msg .. "Cursor: "
    if inv:GetItem(Slot.Cursor).valid then
        msg = string.format("%s%s (%d)\n", msg, inv:GetItem(Slot.Cursor):GetName(), inv:GetItem(Slot.Cursor):GetID())
    else
        msg = string.format("%snone\n", msg)
    end

    msg = msg .. "Target: "
    if client:GetTarget().valid then
        msg = string.format("%s%s (%d)\n", msg, client:GetTarget():GetName(), client:GetTarget():GetID())
    else
        msg = string.format("%snone\n", msg)
    end

    msg = msg .. string.format("Location: `#zone %s %d %d %d`\n", eq.get_zone_id(),  client:GetX() , client:GetY() , client:GetZ())


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
    
    local isClosed, _, code = w:close();
    if not isClosed then
        e.self:Message(MT.White, "Error closing issue file: " .. code);
        return
    end

    e.self:Message(MT.White, "Issue created successfully. Thanks for the report!");
end

return issue;