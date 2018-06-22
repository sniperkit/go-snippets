#!env bash
TERRAFORM_VERSION=${TERRAFORM_VERSION:-0.6.15}
PACKER_VERSION=${PACKER_VERSION:-0.10.0}
cat << EOF
VERSIONS SPECIFIED:
TERRAFORM_VERSION: $TERRAFORM_VERSION
PACKER_VERSION:    $PACKER_VERSION
EOF

terraform_to_keep="
terraform
terraform-provider-aws
terraform-provider-consul
terraform-provider-docker
terraform-provider-github
terraform-provider-null
terraform-provider-packet
terraform-provider-template
terraform-provider-terraform
terraform-provider-tls
terraform-provisioner-file
terraform-provisioner-local-exec
terraform-provisioner-remote-exec
"

# ... hashicorp uris
HASHI_URI="https://releases.hashicorp.com"
TF_BASE_URI="$HASHI_URI/terraform" 
P_BASE_URI="$HASHI_URI/packer" 
HASHIVIM_URI="https://github.com/hashivim"
TF_VIM_URI="$HASHIVIM_URI/vim-terraform.git" 
P_VIM_URI="$HASHIVIM_URI/vim-packer.git" 

#     ... constructed download uris
TF_DL_URI="$TF_BASE_URI/$TERRAFORM_VERSION/terraform_${TERRAFORM_VERSION}_linux_amd64.zip"
P_DL_URI="$P_BASE_URI/$PACKER_VERSION/packer_${PACKER_VERSION}_linux_amd64.zip"

# ... install dirs
BIN_DIR="/usr/local/bin"
TF_ZIP="$BIN_DIR/terraform.zip" 
P_ZIP="$BIN_DIR/packer.zip" 
VIM_DIR="/root/.vim/bundle"

echo "... downloading:"
echo "    $TF_DL_URI -> $TF_ZIP"
wget -q -O $TF_ZIP $TF_DL_URI
echo "    $P_DL_URI -> $P_ZIP"
wget -q -O $P_ZIP $P_DL_URI

unzip $TF_ZIP -d $BIN_DIR
unzip $P_ZIP -d $BIN_DIR

terraform_bins=$(ls -1 $BIN_DIR/terraform* | xargs -n1 -i{} basename {} | sort | uniq)
terraform_to_keep=$(echo "$terraform_to_keep" | sort | uniq | sed '/^$/d')
delete_these=$( comm -13 <(echo "$terraform_to_keep") <(echo "$terraform_bins") | sed -e "s#^#$BIN_DIR/#" )

echo "... will also delete:"
echo "$delete_these"

THIS_DIR=$(pwd)
cd $VIM_DIR
git clone $TF_VIM_URI
git clone $P_VIM_URI
rm -rf $TF_ZIP \
       $P_ZIP \
       $VIM_DIR/terraform/.git \
       $VIM_DIR/packer/.git \
       $delete_these

cd $THIS_DIR


