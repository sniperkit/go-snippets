| Path | Action |
|------|--------|
| /self::\*\[not\(@InternalType='FILE'\)\] | Error |
| /self::\*\[@InternalType='FILE'\] | File |
| /self::\*\[@InternalType='FILE'\]//\*\[@InternalType='\[Bash\] Comment'\] | Comment |
| /self::\*\[@InternalType='FILE'\]//\*\[@InternalType='\[Bash\] shebang element'\] | Comment, Documentation |
| /self::\*\[@InternalType='FILE'\]//\*\[@InternalType='var\-def\-element'\]/\*\[@InternalType='\[Bash\] assignment\_word'\] | SimpleIdentifier |
| /self::\*\[@InternalType='FILE'\]//\*\[@InternalType='function\-def\-element'\]/\*\[@InternalType='\[Bash\] function'\] | FunctionDeclaration |
| /self::\*\[@InternalType='FILE'\]//\*\[@InternalType='function\-def\-element'\]/\*\[@InternalType='\[Bash\] named symbol'\] | FunctionDeclarationName |
| /self::\*\[@InternalType='FILE'\]//\*\[@InternalType='function\-def\-element'\]/\*\[@InternalType='group element'\] | FunctionDeclarationBody, Block |
| /self::\*\[@InternalType='FILE'\]//\*\[@InternalType='if shellcommand'\] | If, Statement |
| /self::\*\[@InternalType='FILE'\]//\*\[@InternalType='for shellcommand'\] | ForEach, Statement |
| /self::\*\[@InternalType='FILE'\]//\*\[@InternalType='while loop'\] | While, Statement |
| /self::\*\[@InternalType='FILE'\]//\*\[@InternalType='until loop'\] | While, Statement |
