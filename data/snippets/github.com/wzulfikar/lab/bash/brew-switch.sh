# switch brew package. 
# eg, switch from php71 to php56: 
# `brew-switch php71 php56`
brew-switch () {
    eval "brew unlink $1 && brew link $2"
}
