#!/bin/bash
# $Id: promoteToEclipse.sh,v 1.3 2007/02/22 08:39:51 asobolev Exp $

# TODO: rename $branch to $version, and $cvsbranch to $branch

# New order of steps in promote (others before or after)
#  0. uploads zips, jars, docs, etc.
#  1. publish RSS
#  2. run parsecvs.sh remotely from   http://build.eclipse.org/modeling/build/updateSearchCVS.php?projects[]=org.eclipse...&projects[]=org.eclipse...
#  3ai. block until done by checking  http://build.eclipse.org/modeling/emf/news/checkReleaseExists.php?project=emf&version={buildIDorAlias}      until value = 1
#  3aii. block until done by checking http://build.eclipse.org/modeling/mdt/news/checkReleaseExists.php?project=uml2-uml&version={buildIDorAlias} until value = 1
#  3b. timeout if stuck for more than x mins
#  4. update bugs (ASSIGNED -> FIXED), announce in newsgroup, etc.

# promoteToEclipse.sh script to update a given build, plus required files on \$downloadServerFullName from CVS
# run script on $buildServerFullName to push data to remote eclipse.org server
# Copyright \(c\) 2004-2006, IBM. Nick Boldt. codeslave\(at\)ca.ibm.com

# requires accompanying properties file, promoteToEclipse.properties, plus a few minimal commandline switches, like -sub ocl -branch 1.0.0 -buildID ...
# can also specify any other properties file with -f flag

# This script references the following other scripts:
# compareFolders (verify upload succeeded)
# fixJavadocs (javadocs)

# see also ../private_html/$projectName/promo.php for web UI to kick this script

###### ###### ###### ###### ###### ###### INSTRUCTIONS & OPTIONS ###### ###### ###### ###### ###### ###### 

norm="\033[0;39m";
grey="\033[1;30m";
green="\033[1;32m";
brown="\033[0;33m";
yellow="\033[1;33m";
blue="\033[1;34m";
cyan="\033[1;36m";
red="\033[1;31m";

function basichelp ()
{
	echo \$Id: promoteToEclipse.sh,v 1.3 2007/02/22 08:39:51 asobolev Exp $
	echo "";
	echo -e "You must specify the properties file used by this script BEFORE any overrides, such as $red-branch$norm or $red-buildID$norm, on the commandline."
	echo "";
	echo -e "usage: $blue$0$norm";
	echo -e "$red-sub$norm          <REQUIRED: MUST BE FIRST FLAG; if ./promoteToEclipse.${red}sub$norm.properties not found, try ./promoteToEclipse.properties>"
	echo -e "$red-branch$norm       <REQUIRED: release version of the files to be promoted, eg., 1.0.0, 1.0.1 (overrides property file)>"
	echo -e "$red-buildID$norm      <REQUIRED: ID of the build>"
	echo "";
	echo -e "$red-cvsbranch$norm    <REQUIRED for RSS feeds: cvs branch of the release, eg., R2_2_maintenance, HEAD>"
	echo -e "$brown-branchIES$norm    <REQUIRED for IES map file: eg., R3_2_maintenance, HEAD>"
	echo ""
	echo -e "$yellow-q$norm, $yellow-Q$norm                   <OPTIONAL: scp, unzip, cvs checkout: quiet or VERY QUIET output>"
	echo -e "$yellow-announce$norm, $yellow-announceonly$norm <OPTIONAL: post announcement in newsgroup - see properties file for settings>"
	echo -e "$yellow-email$norm                   <OPTIONAL: email to notify when done>"
}

function advhelp()
{
	echo "";
	echo -e "$brown-h$norm, $brown-help$norm     <advanced help options>";
}

function morehelp ()
{
	echo ""
	echo -e "-f            <OPTIONAL: if used, MUST BE FIRST FLAG; use specific file instead of ./promoteToEclipse.${red}sub$norm.properties>"
	echo "-user         <username on *.eclipse.org (default is $USER (\$USER))>"
	echo "-debug        <default 0; can be increased to 1 or 2>"
	echo "-tmpfolder    <folder to use when working> (optional)"
	echo "-noclean      <don't clean up when done (leave tmpfolder)>"

	# no/only options (9): drop, jars, rss, IES, docs, searchcvs
	echo ""
	echo "-nodrop, -droponly <do NOT / ONLY upload drop>"
	echo "  -noCompareDropsFolders <after uploading the drop, DO NOT compare source and target for matching MD5s, etc.>"
	echo "-nojars, -jarsonly <do NOT / ONLY do Update Manager jars create & upload step>";
	echo "  -basebuilderBranch  <org.eclipse.releng.basebuilder CVS branch, such as M2_32 or M3_32>"
	echo "  -noCompareUMFolders <after uploading the UM jars, DO NOT compare source and target for matching MD5s, etc.>"
	
	echo ""
	echo "-norss,  -rssonly  <do NOT / ONLY gen RSS feed file update for the given buildID & branch> (optional)"
	echo -e "  $red-cvsbranch$norm       <REQUIRED for RSS feeds: cvs branch of the release, eg., R2_2_maintenance, HEAD>"
	echo "  -feedURL         <override value set in feedPublish.*.properties, eg. for use with N builds>"
	echo "  -feedFile        <override value set in feedPublish.*.properties, eg. for use with N builds>"
	echo ""

	echo "-noIES,  -IESonly  <do NOT / ONLY gen IES map file update for the given buildID & branch> (optional)"
	echo -e "  $brown-branchIES$norm       <REQUIRED for IES map file: eg., R3_2_maintenance, HEAD>"
	echo "  -userIES         <override value set in .properties file, eg., -userIES nboldt>"
	echo ""
	echo "-nodocs, -docsonly           <do NOT / ONLY do javadocs step>";
	echo "-nosearchcvs, -searchcvsonly <do NOT / ONLY do Search CVS database update>";

	echo ""
	echo -e "$grey-nomaven$norm, $grey-mavenonly$norm         <do NOT / ONLY gen Maven jars> (optional)";
	echo -e "$grey-coordsite$norm                   <update Callisto/Europa coordinated release UM site>";
}

function exampleshelp ()
{
	echo ""
	echo -e "To promote a drop, you need at least $red-sub$norm, $red-branch$norm, $red-buildID$norm, $red-cvsbranch$norm, and $brown-branchIES$norm (if applicable), with optional $yellow-announce$norm and $yellow-Q$norm:"
	echo ""
	echo -e "  $blue./promoteToEclipse.sh $red-sub ocl  $yellow-Q -announce $red-branch 1.0.2 $red-buildID M200612071439 $red-cvsbranch R1_0_maintenance $brown-branchIES R3_2_maintenance$norm"
	echo -e "  $blue./promoteToEclipse.sh $red-sub uml2 $yellow-Q -announce $red-branch 2.1.0 $red-buildID I200510101000 $red-cvsbranch HEAD $brown-branchIES HEAD$norm"
	echo ""
	echo "To announce only: ";
	echo -e "  $blue./promoteToEclipse.sh $red-sub emf  $yellow-Q -announceonly $red-branch 2.3.0 -buildID I200611030200$green 2>&1 | tee ~/promo_\`date +%Y%m%d_%H%M%S\`.txt$norm"
	echo ""
}

if [ $# -lt 1 ]; then
	basichelp;
	advhelp;
	exampleshelp;
	exit 1;
elif [[ $1 = "-h" ]] || [[ $1 = "--help" ]]; then
	basichelp;
	morehelp;
	exampleshelp;
	exit 1;
fi

###### ###### ###### ###### ###### ###### BEGIN SETUP ###### ###### ###### ###### ###### ###### 

# default to default properties file
defaultPropertiesFile=./promoteToEclipse.properties
propertiesFiles="";

# for use with -*only flags, zero out all variables
# no/only options (9): drop, jars, rss, IES, docs, searchcvs
# yes/only options(2): (coordsite), announce
function allZero()
{
   dodrop=0;
   UMjars=0;
      RSS=0;
      IES=0;
   dodocs=0;
searchCVS=0;
}

# Create local variable based on the input
echo " "
echo "[promote] Started `date +%H:%M:%S`. Executing with the following options:"
while [ "$#" -gt 0 ]; do
	case $1 in
		'-f')    propertiesFile=$2; 								echo "   $1 [using $propertiesFile]"; if [ -r $propertiesFile ]; then source $propertiesFile; else echo "[promote] Properties file $propertiesFile not found. Exiting..."; exit 99; fi; shift 1;;
		'-sub')  subprojectName=$2; 								echo "   $1 $2";
			# chain them together in order of priority: subproj specific one, default
			propertiesFiles=$propertiesFiles" ./promoteToEclipse."$subprojectName".properties "$defaultPropertiesFile; 
			loaded=0; for propertiesFile in $propertiesFiles; do if [ "$loaded" -eq 0 ] && [ -r $propertiesFile ]; then echo -n "    [loading $propertiesFile ... "; . $propertiesFile; echo "done]"; loaded=1; fi; done
			if [ "$loaded" -eq 0 ]; then echo "    [Can't load any of: $propertiesFiles. Exiting!]";exit 99; fi
			shift 1;;

	# build details

		'-branch')   branch=$2; echo "   $1 $2"; shift 1;;
		'-buildID') buildID=$2; echo "   $1 $2"; shift 1;;

	# user/email

		'-user')     user=$2;   echo "   $1 $2"; shift 1;;
		'-email')   	 email=$2; echo "   $1 $2"; shift 1;;

	# no/only options (9): drop, jars, rss, IES, docs, searchcvs
	
		'-droponly') allZero; dodrop=1; echo "   $1"; shift 0;;
		'-nodrop') 			  dodrop=0; echo "   $1"; shift 0;;
		'-noCompareDropsFolders') noCompareDropsFolders=1; echo "   $1"; shift 0;;

		'-jarsonly') allZero;  UMjars=1; echo "   $1"; shift 0;;
		'-nojars')             UMjars=0; echo "   $1"; shift 0;;
		'-basebuilderBranch') basebuilderBranch="-basebuilderBranch $2"; echo "   $1 $2"; shift 1;;
		'-noCompareUMFolders') noCompareUMFolders="-noCompareUMFolders"; echo "   $1";    shift 0;;

		'-rssonly') allZero; RSS=1;	echo "   $1"; shift 0;;
		'-norss') 			 RSS=0; echo "   $1"; shift 0;;
		'-cvsbranch') cvsbranch=$2;	echo "   $1 $2"; shift 1;;
		'-feedURL')	  feedURL="-DfeedURL=$2"; echo "   $1 $2"; shift 1;;
		'-feedFile')  feedFile="-Dfile=$2";   echo "   $1 $2"; shift 1;;

		'-IESonly') allZero; IES=1;	echo "   $1"; shift 0;;
		'-noIES') 			 IES=0; echo "   $1"; shift 0;;
		'-branchIES') branchIES=$2;	echo "   $1 $2"; shift 1;;
		'-userIES') userIES=$2; 	echo "   $1 $2"; shift 1;;

		'-docsonly') allZero; dodocs=1; echo "   $1"; shift 0;;
		'-nodocs') 			  dodocs=0; echo "   $1"; shift 0;; 

		'-searchcvsonly') allZero; searchCVS=1; echo "   $1"; shift 0;;
		'-nosearchcvs') 		   searchCVS=0; echo "   $1"; shift 0;; 

	# yes/only options(2): (coordsite), announce
	
		'-coordsite') 			   coordsite=1; echo "   $1"; shift 0;;
		'-coordsiteonly') allZero; coordsite=1; echo "   $1"; shift 0;;

		'-announce') 			   announce=1;  echo "   $1"; shift 0;;
		'-announceonly') allZero;  announce=1;  echo "   $1"; shift 0;;

	# debug options
	
		'-tmpfolder') tempfold=$2; echo "   $1 $2";	shift 1;;
		'-q') quietCVS=-q; quiet=-q;  echo "   $1"; shift 0;;
		'-Q') quietCVS=-Q; quiet=-q;  echo "   $1"; shift 0;;
		'-debug') debug=$2; echo "   $1 $2"; shift 1;;
		'-noclean') noclean=1; echo "   $1"; shift 0;;
	esac
	shift 1
done

if [ "$subprojectName" = "" ]; then # no value set!
  echo "[promote] No subproject name set in properties file or by -sub flag. Script cannot continue. Exiting...";
  exit 99;
fi

if [ "$branch" = "" ]; then # no value set!
  echo "[promote] No branch value set in properties file or by -branch flag. Script cannot continue. Exiting...";
  exit 99;
fi

if [ "$buildID" = "" ]; then # no value set!
  echo "[promote] No build ID value set in properties file or by -buildID flag. Script cannot continue. Exiting...";
  exit 99;
fi

# get path to PHP interpreter
PHP=php
if [ -x /usr/bin/php ]; then
	PHP=/usr/bin/php
elif [ -x /usr/bin/php4 ]; then
	PHP=/usr/bin/php4
elif [ -x /usr/bin/php5 ]; then
	PHP=/usr/bin/php5
else
	PHP=php
fi

# create default temp folder
mkdir -p /home/$user/tmp;

###### ###### ###### ###### ###### VARIABLES DERIVED FROM $user ###### ###### ###### ###### ###### 

# temp folder base (make 'em unique so concurrent builds don't overlap)
tempfold=/home/$user/tmp/promoteToEclipse-$projectName-$subprojectName-$user-`date +%Y%m%d_%H%M%S`

#users (for ssh and cvs connections)
buildServerCVSUser=$user"@"$buildServerFullName
eclipseCVSUser=$user"@"$eclipseServerFullName
eclipseSSHUser=$user"@"$downloadServerFullName
eclipseSSHUserHome=$(ssh $eclipseSSHUser "echo \$HOME");

#cvs paths
buildServerCVSRep=:ext:$buildServerCVSUser:/home/cvs
wwwCVSRep=:ext:$eclipseCVSUser:/cvsroot/org.eclipse
eclipseCVSRep=:ext:$eclipseCVSUser:/cvsroot/technology
anonEclipseCVSRep=:pserver:anonymous@dev.eclipse.org:/cvsroot/eclipse
coordsiteCVSRep=:ext:$eclipseCVSUser:/cvsroot/callisto

# if no value for userIES, default to commandline -user or value in properties file
if [ "$userIES" = "" ]; then userIES=$user; fi
IESCVSRep=":ext:"$userIES"@ottcvs1.ottawa.ibm.com:/home/cvs/com.ibm.ies.releng"

###### ###### ###### ###### ###### SETUP DONE, DEFINE METHODS ###### ###### ###### ###### ###### 

checkZipExists ()
{
	theFile=`echo $1 | sed -e 's/^.*\///'`
	theURL=$1
	theDir=$2
	$ANT -f $buildScriptsDir/checkZipExists.xml getZip -DdownloadsDir=$theDir -DtheFile=$theFile -DtheURL=$theURL
	#echo "[start] Ant returned: $#"
}

getBuildIDactual ()
{
	#new, more efficient method as of nov 12 thanks to ken's identification of the old way's limitation
	buildIDactual=`find $buildDropsDir/$branch/$buildID -name "$SDKfilenamepattern"`
	buildIDactual=${buildIDactual##*SDK-}; # trim up to SDK- (substring notation)
	buildIDactual=${buildIDactual%%\.zip}; # trim off .zip (substring notation)
	#echo $buildIDactual
}

sendEmail ()
{
	if [ "x$email" != "x" ]; then
		ssh $buildServerCVSUser "$PHP -q $buildScriptsDir/sendEmail.php -email $email -projectName $projectName -branch $branch -buildID $buildID -promote true";
	fi
}

###### ###### ###### ###### METHODS DONE, BEGIN WORK HERE ###### ###### ###### ###### 

buildIDactual=buildID;

if [ $dodrop -eq 0 ]; then
	echo "[promote] Upload new drop - omitted."
else
	echo "[promote] Remove any temp files left over from a -noclean build"
	for f in \
		`find $buildDropsDir/$branch/$buildID -type d -name "org.eclipse*releng*"` \
		$buildDropsDir/$branch/$buildID/testing $buildDropsDir/$branch/$buildID/eclipse; do
		if [ -d $f ]; then rm -fr $f; fi
	done		
	echo "[promote] [`date +%H:%M:%S`] Create target drop folder $projectDropsDir/$branch/$buildID on remote box started:"
	ssh $eclipseSSHUser mkdir -p $projectDropsDir/$branch/$buildID
	echo "[promote] [`date +%H:%M:%S`] Create target drop folder $projectDropsDir/$branch/$buildID on remote box done."

	echo "[promote] [`date +%H:%M:%S`] SCP build folder $branch/$buildID onto remote box started:"
	scp -r -v $buildDropsDir/$branch/$buildID $eclipseSSHUser:$projectDropsDir/$branch/
	echo "[promote] [`date +%H:%M:%S`] SCP build folder $branch/$buildID onto remote box done."

	if [ $noCompareDropsFolders -eq 0 ]; then
		### CHECK MD5s and compare dir filesizes for match (du -s)
		echo "[promote] [`date +%H:%M:%S`] Comparing local and remote folders to ensure SCP completeness... "
		$buildScriptsDir/compareFolders.sh -user $user -local $buildDropsDir/$branch/$buildID -remote $projectDropsDir/$branch/$buildID -server $eclipseSSHUser
		returnCode=$?
		if [ $returnCode -gt 0 ]; then
			echo "[promote] [`date +%H:%M:%S`] ERROR! Script exiting with code $returnCode from compareFolders.sh"
			exit $returnCode;
		fi
	else
		echo "[promote] [`date +%H:%M:%S`] Comparing local and remote folders to ensure SCP completeness ... omitted."
	fi

	# new 040429/5pm
	echo "[promote] [`date +%H:%M:%S`] CHMOD build folder $branch/$buildID on remote box to give group perms started:"
		ssh $eclipseSSHUser "
			chmod -R $eclipsePermsDir $projectDropsDir/$branch/$buildID
			chgrp -fR $eclipseUserGroup $projectDropsDir/$branch/$buildID
		"
		## must keep closing quote on preceeding line to close the ssh command section
	echo "[promote] [`date +%H:%M:%S`] CHMOD build folder $branch/$buildID on remote box to give group perms done."
fi

if [ $UMjars -eq 0 ]; then
	echo "[promote] Create & promote Update Manager jars to download ... omitted [-nojars]."
else
	echo "[promote] [`date +%H:%M:%S`] Create & promote Update Manager jars to download started:"
	echo "[promote] Running buildUpdate.sh:"
	ssh $buildServerCVSUser "cd $buildScriptsDir; ./buildUpdate.sh -sub "$subprojectName" -user $user $quietCVS -branch $branch -buildID $buildID -promote $noCompareUMFolders $basebuilderBranch -debug $debug";
	echo "[promote] [`date +%H:%M:%S`] Create & promote Update Manager jars to download done."
fi

	### Update Eclipse Project RSS feed file
if [ "$RSS" -eq 0 ]; then
	echo "[promote] Update RSS feed - omitted."
elif [ "$RSS" -eq 1 ]; then
	echo "[promote] Update RSS feed ..."

	cd $buildScriptsDir;

	if [ "$cvsbranch" = "" ]; then # no value set!
	  echo "[promote] No CVS branch value set in properties file or by -cvsbranch flag. Script cannot continue. Exiting...";
	  exit 99;
	fi
		
	getBuildIDactual; buildAlias=$buildIDactual;
	if [[ $debug -gt 0 ]]; then echo "[promote] Using buildAlias = $buildIDactual, branch = $branch, buildID = $buildID"; fi

	# find dependencyURLs from $buildDropsDir/$branch/$buildID/buildlog.txt
	# http://download.eclipse.org/downloads/drops/M20060609-1217/eclipse-SDK-M20060609-1217-linux-gtk.tar.gz
	# http://download.eclipse.org/tools/emf/downloads/drops/2.2.0/I200606150000/emf-sdo-xsd-SDK-I200606150000.zip
	dependencyURLs="";
	for dep in `head -50 $buildDropsDir/$branch/$buildID/buildlog.txt | grep "\-URL" | sed -e "s/^[^hf]\\+//g" -e "s/fullmoon\..\\+\.ibm\.com/download.eclipse.org/g" | sort | uniq`; do 
		if [ "x$dependencyURLs" != "x" ]; then dependencyURLs="$dependencyURLs,"; fi # join with commas
		dependencyURLs=$dependencyURLs"$dep";
		if [[ $debug -gt 0 ]]; then echo "[promote] Using URL = $dep"; fi			
	done
	# http://download.eclipse.org/downloads/drops/M20060609-1217/eclipse-SDK-M20060609-1217-linux-gtk.tar.gz,http://download.eclipse.org/tools/emf/downloads/drops/2.2.0/I200606150000/emf-sdo-xsd-SDK-I200606150000.zip
	
	export JAVA_HOME=/opt/sun-java2-5.0;
	export ANT_HOME=/home/www-data/apache-ant-1.6.5;

	# extract jar from CVS if not available (first time only - this will only work if run as www-data due to write permissions)
	if [ ! -f feedTools.jar ]; then $ANT_HOME/bin/ant -f feedWatchSetup.xml; fi

	buildType=${buildID:0:1};
	typeSuffix=""; if [ "$buildType" = "N" ]; then typeSuffix="\-N"; fi
	
	mkdir -p $tempfold/rss; 
	# replace %%project%% (ocl, transaction, etc.); %%typeSuffix%% ("-N" or ""); %%user%% (cdamus, nickb, etc.) with actual values
	cat feedPublish.properties | sed -e "s/%%user%%/$user/g" -e "s/%%typeSuffix%%/$typeSuffix/g" -e "s/%%project%%/$subprojectName/g" > $tempfold/rss/feedPublish.$subprojectName.properties
	
	CLASSPATH="$JAVA_HOME/lib/rt.jar:"$ANT_HOME/lib/ant.jar":"$ANT_HOME/lib/ant-launcher.jar;
	cmd="$JAVA_HOME/bin/java -Dant.home=$ANT_HOME -Dant.library.dir=$JAVA_HOME/lib -classpath $CLASSPATH org.apache.tools.ant.launch.Launcher"
	cmd=$cmd" -buildfile feedPublish.xml -propertyfile $tempfold/rss/feedPublish.$subprojectName.properties"
	
	# override values in properties file with current values (or assign dynamic/missing values)
	# TODO: rename $branch to $version, and $cvsbranch to $branch
	cmd=$cmd" -Dproject=$subprojectName -Dversion=$branch -Dbranch=$cvsbranch -DbuildID=$buildID -DbuildAlias=$buildAlias -DbuildType=$buildType";
	cmd=$cmd" -DdependencyURLs=$dependencyURLs -Ddebug=$debug $feedURL $feedFile";
	if [ "$buildType" == "N" ]; then # override default buildURL since these are on build server, not published
		cmd=$cmd" -DbuildURL=http://emf.torolab.ibm.com/tools/emf/downloads/drops/"$branch"/"$buildID;
	fi
	echo ""; echo $cmd | perl -pe "s/ -/\n  -/g" | perl -pe "s/\.jar:/\.jar\n    :/g"; echo "";
	$cmd;
	
	# note: created files will not have the correct group ownership / perms
	echo "[promote] feed update done."
fi

if [ "$IES" -eq 0 ]; then
	echo "[promote] Update IES map file - omitted."
elif [ "$IES" -eq 1 ]; then
	# update the IES mapping file: checkout file, edit it, then check back in (commit)
	echo "[promote] Update IES map file com.ibm.ies.releng/maps/$projectName-$subprojectName.map"

	if [ "$branchIES" = "" ]; then # no value set!
  		echo "[promote] No IES branch value set in properties file or by -branchIES flag. Script cannot continue. Exiting...";
  		exit 99;
	fi

	getBuildIDactual;
	# create file: /home/cvs/com.ibm.ies.releng/com.ibm.ies.releng/maps/$projectName-$subprojectName.map
	IESmapfile=$projectName"-"$subprojectName".map";
	echo "[promote] Got actual buildID: $buildIDactual, IES branch: $branchIES, mapfile: $IESmapfile";

	# setup
	tmpfolder=$tempfold/ies; mkdir -p $tmpfolder/1; cd $tmpfolder;

	# checkout
	if [ "$branchIES" != "HEAD" ]; then branchIES="-r "$branchIES; else branchIES=""; fi
	cvscmd="cvs -d $IESCVSRep $quietCVS co $branchIES -P -d 1 com.ibm.ies.releng/maps";
	if [[ $debug -gt 0 ]]; then echo "[promote] [`date +%H:%M:%S`] "$cvscmd; fi
	$cvscmd;
	echo "[promote] [`date +%H:%M:%S`] done."

	webPath=http://fullmoon.ottawa.ibm.com/technology/$projectName/$subprojectName/downloads/drops/$branch/$buildID/

	#edit file, then make a copy in the other folder
	echo "[promote] Writing $IESmapfile ...";
	outfile=$tmpfolder/1/$IESmapfile

	index=0;
	element_count=${#IESmapfileArray[@]}
	while [ "$index"  -lt "$element_count" ]; do
	  # eg., echo "$projectName-SDK-$buildIDactual.zip=$webPath     | | | sdk | $projectName"    >> $outfile.tmp ;
	  txt=${IESmapfileArray[$index]/buildIDactual/$buildIDactual}; txt=${txt/webPath/$webPath};
	  echo $txt >> $outfile.tmp;
	  let "index = $index + 1";
	  let index = index + 1;
	done

	## adds optional cvs add command
	cd $tmpfolder/1 ;
	if [ ! -f $IESmapfile ]; then
		echo "[promote] [`date +%H:%M:%S`] add $IESmapfile";	## add if not found
		mv -f $IESmapfile.tmp $IESmapfile ;
		cvs -d $IESCVSRep $quietCVS add -k v -m "promoteToEclipse: new map file" $IESmapfile;
	else
		dif=`diff -q $IESmapfile.tmp $IESmapfile`; # echo $dif;
		if [ "x$dif" = "x" ]; then # remove temp file, it's the same as the previous
			echo "[promote] [`date +%H:%M:%S`] $IESmapfile exists, no change ";	## no add if found
			rm -f $IESmapfile.tmp ;
		else # files differ
			echo "[promote] [`date +%H:%M:%S`] $IESmapfile exists, updating ";	## no add if found
			mv -f $IESmapfile.tmp $IESmapfile ;
		fi
	fi

	#check in
	cvs -d $IESCVSRep $quietCVS ci -m "promoteToEclipse: $branch/$buildID";
	echo "[promote] [`date +%H:%M:%S`] done." ;
	
fi

if [ "$dodocs" -eq 0 ]; then
	echo "[promote] Unzip new javadocs - omitted."
else

	echo "[promote] [`date +%H:%M:%S`] Javadoc creation started: "
	index=0;
	element_count=${#javadocTargetArray[@]}
	while [ "$index"  -lt "$element_count" ]; do
	  ssh $buildServerCVSUser "cd $buildScriptsDir; ./fixJavadocs.sh "${javadocTargetArray[$index]}" $quiet -branch $branch -buildID $buildID";
	  let "index = $index + 1"
	done
	echo "[promote] [`date +%H:%M:%S`] Javadoc creation done."
fi

#  run parsecvs.sh remotely     http://build.eclipse.org/modeling/build/updateSearchCVS.php?projects[]=org.eclipse...&projects[]=org.eclipse...
#  block until done by checking http://build.eclipse.org/modeling/emf/news/checkReleaseExists.php?project=emf&version={buildIDorAlias}      until value = 1
#  block until done by checking http://build.eclipse.org/modeling/mdt/news/checkReleaseExists.php?project=uml2-uml&version={buildIDorAlias} until value = 1
#  timeout if stuck for more than x mins

# TODO: migrate checks and updates to build.eclipse.org (must update /modeling/includes/searchcvs-dbaccess.php in /www/ root)
searchCVSServer="http://emft.eclipse.org"

if [[ $searchCVS -eq 1 ]]; then
  echo "[promote] [`date +%H:%M:%S`] Update Search CVS & Release Notes database ...";
  parsecvsCmd="$searchCVSServer/modeling/build/updateSearchCVS.php?";
  parsecvsCmd=$parsecvsCmd"projects%5B%5D="org.eclipse.emft"&";
  tmpfile=`mktemp`;
  cmd="wget -q -O $tmpfile \"$parsecvsCmd\"";
  echo "[promote] $cmd";
  $cmd;
  rm -f $tmpfile;
else
  echo "[promote] [`date +%H:%M:%S`] Update Release Notes & Search CVS database ... omitted.";
fi

getBranchLabel () 
{
		if [ $buildIDactual == $branch ]; then						# echo "branch and actual are equal, need only $branch";
			branchLabel="$branch";
		elif [ ${buildIDactual##$branch} != $buildIDactual ]; then	# echo "just need $buildIDactual";
			branchLabel="$buildIDactual";
		else														# echo "Need both $branch and $buildIDactual";
			branchLabel="$branch $buildIDactual";
		fi
}

if [[ $coordsite -eq 1 ]]; then
	if [[ ! $coordsiteName ]];   then coordsiteName="europa"; fi # set default if missing
	if [[ ! $coordsiteBranch ]]; then coordsiteBranch="HEAD"; fi # set default if missing
	
	# setup
	tmpfolder=$tempfold/coordsite; mkdir -p $tmpfolder; cd $tempfold;
	qualifier=".v"${buildID:1};
	buildType=${buildID:0:1};
	if [ $buildType = "R" ]; then sitexml="site.xml"; else sitexml="site-interim.xml"; fi
	
	featuresXML="features-$projectName-$subprojectName.xml";
	
	echo "[coordsite] [`date +%H:%M:%S`] Update Coordinated Update Site ($sitexml):";
	# checkout
	if [ "$coordsiteBranch" != "HEAD" ]; then coordsiteBranch="-r "$coordsiteBranch; else coordsiteBranch=""; fi
	cvscmd="cvs -d $coordsiteCVSRep $quietCVS co $coordsiteBranch -d coordsite org.eclipse."$coordsiteName".tools/build-home";
	if [[ $debug -gt 0 ]]; then echo "[coordsite] [`date +%H:%M:%S`] "$cvscmd; fi
	$cvscmd;
	
	# get actual feature versions for the following defined list: org.eclipse.*sub_x.y.z.vqualifier.jar
	#TODO may not work for all subprojects, eg. jeteditor, transaction (workspace)
	cd $localWebDir/updates/features; features="`find . -regex \".*/org\.eclipse\..*$subprojectName_.*$qualifier\.jar\"`";

	cd $tmpfolder;
	echo -n '' > $featuresXML
	echo '<?xml version="1.0" encoding="UTF-8"?>
<project name="update" default="update">
    <target name="update">

      <property name="updateSite"
                value="file://${user.home}/downloads/technology/emft/updates/'$sitexml'" />

        <echo message="   pulling update jars from ${updateSite}" />' >> $featuresXML
	for feat in $features; do
		f=${feat:2}; # trim first two chars
		featureId=${f%_*}; # trim from _ to end
		featureVersion=${f#*_}; featureVersion=${featureVersion%.v*}; # trim from start to _, then from .v to end
		echo "            "$featureId" "$featureVersion""$qualifier;
echo '
        <ant antfile="updateMirrorProject.xml">
            <property name="featureId"
                      value="'$featureId'" />
            <property name="version"
                      value="'$featureVersion''$qualifier'" />
        </ant>' >> $featuresXML;
	done
echo '
    </target>
</project>
' >> $featuresXML;
	#cat $featuresXML;
	
	# commit changes
	cvs ci -m "promoteToEclipse: $sitexml, $branch$qualifier" $featuresXML; 
    #echo "cvs ci -m \"promoteToEclipse: $sitexml, $branch$qualifier\" $featuresXML"; 
	
	echo "[coordsite] [`date +%H:%M:%S`] done.";
fi
	
if [[ $noclean -eq 0 ]]; then
	# cleanup temp space
	ssh $buildServerCVSUser "rm -fr $tempfold/";
	echo "[promote] Temporary files purged from $tempfold/."
else
	echo "[promote] Temporary files left in $tempfold/. Please scrub them manually when done."
fi

# send newsgroup notification
postfile=newsgroup-post.txt
if [[ $announce -eq 1 ]]; then
	echo -n "[announce] Announce new build in $projectName newsgroup... ";
	if [ -r /home/$user/.ssh/$postfile ]; then
		echo "";
		getBuildIDactual;
		getBranchLabel;
		mkdir -p $tempfold;
		head -2 /home/$user/.ssh/$postfile > $tempfold/$postfile;
		# this file should contain the following two lines:
		# authinfo user exquisitus
		# authinfo pass ********** (replace with real pwd)
		echo "post" >> $tempfold/$postfile;
		echo "From: $newsgroupPublisherEmail" >> $tempfold/$postfile;
		echo "Subject: [Announce] $projectNameVanity $branchLabel is available" >> $tempfold/$postfile;
		echo "Newsgroups: $newsgroup" >> $tempfold/$postfile;
		if [ "x$newsgroupThreadReferences" != "x" ]; then # add support for threading under an existing post
			echo "References: $newsgroupThreadReferences" >> $tempfold/$postfile;
			echo "In-Reply-To: $newsgroupThreadReferences" >> $tempfold/$postfile;
		fi
		echo "" >> $tempfold/$postfile;
		echo "$releaseNotesURL?version=$branch" >> $tempfold/$postfile;
		echo "$downloadsURL" >> $tempfold/$postfile;
		echo "" >> $tempfold/$postfile;
		echo "." >> $tempfold/$postfile;
		echo "quit" >> $tempfold/$postfile;
		cat $tempfold/$postfile | nc news.eclipse.org 119 2>&1 ;
		rm -fr $tempfold/$postfile;
		echo "[announce] done.";
	else
		echo "$postfile not found or not readable, message delivery disabled."
	fi
fi

sendEmail;

echo "[promote] Done `date +%H:%M:%S`"
echo ""

