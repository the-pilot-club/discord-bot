package static_text

import (
	"fmt"
	"tpc-discord-bot/internal/config"
)

func SetSOPRolesText(guildId string) string {
	var output string

	output += "📑 **Server Roles:**\n\n"
	output += fmt.Sprintf("<:pilots:895488969524006954> <@&%v> - Sim pilots\n", config.GetRoleId(guildId, "Pilots"))
	output += fmt.Sprintf("<:TPC_Charters:897128289234190376> <@&%v> -  TPC Charters VA (OnAir platform)\n", config.GetRoleId(guildId, "Charters Pilots"))
	output += fmt.Sprintf("<:IRL_pilots:895488931922067517> <@&%v> - Real life pilots\n", config.GetRoleId(guildId, "IRL Pilots"))
	output += fmt.Sprintf("<:atc:895492003528335400> <@&%v> - Air Traffic Controllers\n", config.GetRoleId(guildId, "ATC"))
	output += fmt.Sprintf("<:early_adopter:936005453559771186> <@&%v> - Members who supported the Club early on\n", config.GetRoleId(guildId, "Early Adopters"))
	output += fmt.Sprintf("<:onboarded:967794040592535622> <@&%v> - Members, who completed TPC onboarding\n", config.GetRoleId(guildId, "Onboarded"))
	output += fmt.Sprintf("<:staff:935982118406856744> <@&%v>  - Members of TPC Staff Team\n", config.GetRoleId(guildId, "Staff"))
	output += fmt.Sprintf("<:team_lead:936003597374722058> <@&%v> - Staff members leading specialty teams\n", config.GetRoleId(guildId, "Team Lead"))
	output += fmt.Sprintf("<:flight_ops:895489026038042634> <@&%v> - Flight coordinators\n", config.GetRoleId(guildId, "Flight Ops"))
	output += fmt.Sprintf("<:training:895480894901592074> <@&%v> - Training Team members\n", config.GetRoleId(guildId, "Training Team"))
	output += fmt.Sprintf("<:training:895480894901592074> <@&%v> - Training Team Instructors\n", config.GetRoleId(guildId, "Flight Instructor"))
	output += fmt.Sprintf("<:media_team:895481110132326420> <@&%v> - Media / Socials \n", config.GetRoleId(guildId, "Social Media Team"))
	output += fmt.Sprintf("<:streaming:895480920960819211> <@&%v> - Streamers\n", config.GetRoleId(guildId, "Streamers"))
	output += fmt.Sprintf("<:partners:895481092524609576> <@&%v> - Club Partners\n", config.GetRoleId(guildId, "Partners"))
	output += fmt.Sprintf("<:boosters:896214393472286730> <@&%v> - Server boosters \n", config.GetRoleId(guildId, "Booster"))
	output += fmt.Sprintf("<:lucky_pilots:895489289767497749> <@&%v> - Giveaway, screenshot contest, or event winners\n", config.GetRoleId(guildId, "Lucky Pilots"))
	output += fmt.Sprintf("<:explorer:907581955766382603> <@&%v> - Explorer Missions series participants\n", config.GetRoleId(guildId, "Explorers"))
	output += fmt.Sprintf("<:long_haul:1152604323793080381> <@&%v> - Long Haul Missions series participants\n", config.GetRoleId(guildId, "Long Haul"))
	output += fmt.Sprintf("<:devs:936002576728588328> <@&%v> - Software and web development team\n", config.GetRoleId(guildId, "Tech Team"))
	output += fmt.Sprintf("<:ground_crew:895489218804064287> <@&%v> - Moderators\n", config.GetRoleId(guildId, "Moderator"))
	output += fmt.Sprintf("<:marshalls:895489071097454632> <@&%v> - Admin\n\n", config.GetRoleId(guildId, "Air Marshals"))
	output += "**Leaderboard Roles**\n\n"
	output += fmt.Sprintf("<@&%v> - Level 15\n", config.GetRoleId(guildId, "Commuter"))
	output += fmt.Sprintf("<@&%v> - Level 25\n", config.GetRoleId(guildId, "Frequent Flyer"))
	output += fmt.Sprintf("<@&%v> - Level 35\n\n", config.GetRoleId(guildId, "VIP"))
	output += "**Self Assigned Roles**\n\n"
	output += fmt.Sprintf("<:giveaway:895480872243978280> <@&%v> - Giveaway participants\n", config.GetRoleId(guildId, "Giveaway"))
	output += fmt.Sprintf("<:GA_Gang:898533440532672512> <@&%v> - General Aviation flight announcements\n", config.GetRoleId(guildId, "GA Flights"))
	output += fmt.Sprintf("<:group_flights:938262268972511373> <@&%v> - Group Flight announcements\n", config.GetRoleId(guildId, "Group Flights"))
	output += fmt.Sprintf("<:streaming:895480920960819211> <@&%v> - YT or Twitch streamers\n", config.GetRoleId(guildId, "Streamers"))
	output += fmt.Sprintf("<:Insiders:996214296503652383> <@&%v> - Opt-in to get pinged with Club updates to stay in the loop\n", config.GetRoleId(guildId, "Insiders"))
	output += fmt.Sprintf("<:training:895480894901592074> <@&%v> : TPC Flight School participants\n", config.GetRoleId(guildId, "Flight School"))
	output += fmt.Sprintf("<:helicopter:1167835494869110784> <@&%v>: Helicopter Pilots\n\n", config.GetRoleId(guildId, "Helicopter"))
	output += "Click the button below to get self-assigned roles:\n "

	return output

}
