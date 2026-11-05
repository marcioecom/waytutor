package discord

import "github.com/bwmarrin/discordgo"

func tutorCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "tutor",
		Description: "Administração do WayTutor",
		Options: []*discordgo.ApplicationCommandOption{
			subcommand("pending", "Lista os cadastros aguardando liberação"),
			subcommand("approve", "Libera uma conta", emailOption()),
			subcommand("pause", "Pausa uma conta", emailOption()),
			subcommand("revoke", "Encerra as sessões de uma conta", emailOption()),
		},
	}
}

func subcommand(name, description string, options ...*discordgo.ApplicationCommandOption) *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        name,
		Description: description,
		Options:     options,
	}
}

func emailOption() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        "email",
		Description: "E-mail da conta",
		Required:    true,
		MaxLength:   200,
	}
}
