package medkit

import (
	"github.com/spf13/cobra"
	"github.com/LeadPipeSoftware/medkit/internal/dotfile"
	"github.com/spf13/viper"
)

var alwaysSkip = false
var alwaysOverwrite = false

var installDotfilesCmd = &cobra.Command{
	Use:   "dotfiles",
	Short: "install .symlink files in your dotfiles directory",
	Long: `
This command looks for any file with a .symlink extension in your dotfiles
directory. When it finds a match, it will create a symbolic link from that
file to your home directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		dotfilesDirectory := viper.GetString(DotFilesDirectory)
		homeDirectory := viper.GetString(HomeDirectory)
		backupExtension := viper.GetString(BackupExtension)

		dotfile.InstallDotfiles(dotfilesDirectory, homeDirectory, alwaysSkip, alwaysOverwrite, backupExtension)
	},
}

func init() {
	installDotfilesCmd.Flags().BoolVarP(&alwaysSkip, "always-skip", "s", false, "always skip existing files")
	installDotfilesCmd.Flags().BoolVarP(&alwaysOverwrite, "always-overwrite", "o", false, "always overwrite existing files")
}

