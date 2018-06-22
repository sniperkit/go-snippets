/*
 * Copyright 2014-2015 The DevMine Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package ch.devmine.javaparser;

import ch.devmine.javaparser.structures.Language;
import ch.devmine.javaparser.structures.Package;
import ch.devmine.javaparser.structures.Project;
import ch.devmine.javaparser.structures.SourceFile;
import ch.devmine.javaparser.utils.FileTypeFinder;
import ch.devmine.javaparser.utils.FileWalkerUtils;
import ch.devmine.javaparser.utils.GsonFactory;
import ch.devmine.javaparser.utils.Log;
import com.github.javaparser.ParseException;
import com.google.gson.Gson;
import java.io.File;
import java.io.FileInputStream;
import java.io.FileNotFoundException;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import org.apache.commons.compress.archivers.tar.TarArchiveEntry;
import org.apache.commons.compress.archivers.tar.TarArchiveInputStream;

/**
 * Main class of the parser
 *
 * @author lweingart
 */
public class Main {

    /**
     * class tag
     */
    public static final String TAG = "Main class::";

    /**
     * language
     */
    private static final String JAVA = "java";

    /**
     * paradigms
     */
    private static final String OOP = "object oriented";
    private static final String STRUCTURED = "structured";
    private static final String IMPERATIVE = "imperative";
    private static final String GENERIC = "generic";
    private static final String REFLECTIVE = "reflective";
    private static final String CONCURRENT = "concurrent";

    /**
     * Main function
     * @param args
     *      the path of the repository to parse
     */
    public static void main(String[] args) {

        // parse one repository
        if (args.length != 1) {
            System.err.println("usage : javaparser <path>\n" +
                               "path : the path to the folder or the tar archive to parse");
            return;
        }

        Project project = new Project();
        Language language = defineJavaLang();
        List<Language> languages = new ArrayList<>();
        languages.add(language);
        project.setLanguages(languages);

        HashMap<String, Package> packs =  new HashMap<>();


        if (new File(args[0]).isDirectory()) {
            parseAsDirectory(project, packs, languages, language, args[0]);
        } else {
            parseAsTarArchive(project, packs, languages, language, args[0]);
        }

        List<Package> packages = new ArrayList<>(packs.values());
        int projLoc = 0;
        for (Package pack : packages) {
            int packLoc = 0;
            for (SourceFile file : pack.getSourceFiles()) {
                packLoc += file.getLoc();
            }
            pack.setLoc(packLoc);
            projLoc += packLoc;
        }
        project.setPackages(packages);
        project.setLoc(projLoc);
        Gson gson = GsonFactory.build();
        String jsonProject = gson.toJson(project);

        // the result is written in the system output in order to be
        // used in chain with the source analyzer
        // see https://github.com/devmine/scranlzr
        System.out.println(jsonProject);
    }

    private static void parseAsTarArchive(Project project, HashMap<String, Package> packs,
                                          List<Language> languages, Language language, String arg) {

        int slashIndex = arg.lastIndexOf("/");
        String projectName = arg.substring(slashIndex + 1, arg.lastIndexOf("."));
        String projectRoute = arg.substring(0, slashIndex + 1);
        project.setName(projectName);

        try {
            TarArchiveInputStream tarInput = new TarArchiveInputStream(new FileInputStream(arg));
            TarArchiveEntry entry;
            while (null != (entry = tarInput.getNextTarEntry())) {
                String entryName = entry.getName();
                String fileName = entryName.substring(entryName.lastIndexOf("/") + 1);
                if (entryName.endsWith(".java") && !(fileName.startsWith("~") || fileName.startsWith("."))) {
                    // add package containing source file to hashmap
                    Package pack = new Package();
                    String packName = FileWalkerUtils.extractFolderName(entry.getName());
                    String packPath = extractPackagePath(entry.getName());
                    pack.setName(packName);
                    pack.setPath(packPath);
                    if (packs.get(packName) == null) {
                        packs.put(packName, pack);
                    }

                    // parse java file

                    String entryPath = projectRoute.concat(entryName);
                    Parser parser = new Parser(entryPath, tarInput);
                    SourceFile sourceFile = new SourceFile();
                    parseAndFillSourceFile(parser, sourceFile, entryPath, language, packs, packName);
                }
            }
        } catch (IOException ex) {
            Log.e(TAG, ex.getMessage());
        }

    }

    private static void parseAsDirectory(Project project, HashMap<String, Package> packs,
                                         List<Language> languages, Language language, String arg) {
        File repoRoot = new File(arg);
        Path repoPath = repoRoot.toPath();

        FilesWalker fileWalker = new FilesWalker();
        try {
            Files.walkFileTree(repoPath, fileWalker);
        } catch (IOException ex) {
            Log.e(TAG, ex.getMessage());
        }

        String path = repoPath.toString();
        String name = path.substring(path.lastIndexOf("/") + 1);
        project.setName(name);

        for (String dir : fileWalker.getDirectories()) {
            dir = dir.replaceAll("\n", "");
            if (FileTypeFinder.search(dir, ".java")) {
                Package pack = new Package();
                String packName = dir.substring(dir.lastIndexOf("/") + 1);
                pack.setName(packName);
                pack.setPath(dir);
                packs.put(packName, pack);
            }
        }

        List<String> javaFiles = FileWalkerUtils.extractJavaFiles(fileWalker.getRegularFiles());
        for (String javaFile : javaFiles) {
            try {
                Parser parser = new Parser(javaFile, null);
                SourceFile sourceFile = new SourceFile();
                String packName = FileWalkerUtils.extractFolderName(javaFile);
                parseAndFillSourceFile(parser, sourceFile, javaFile, language, packs, packName);
            } catch (FileNotFoundException ex) {
                Log.e(TAG, ex.getMessage());
            }
        }

    }

    private static void parseAndFillSourceFile(Parser parser, SourceFile sourceFile, String javaFile, Language language,
                                               HashMap<String, Package> packs, String packName) {
        try {
            parser.parse();
            if (!parser.isPackageInfo()) {
                sourceFile.setPath(javaFile);
                sourceFile.setLanguage(language);
                sourceFile.setImports(parser.getImports());
                sourceFile.setInterfaces(parser.getInterfaces());
                sourceFile.setClasses(parser.getClasses());
                sourceFile.setEnums(parser.getEnums());
                sourceFile.setLoc(parser.getLoc());

                packs.get(packName).getSourceFiles().add(sourceFile);
            } else {
                packs.get(packName).setDoc(parser.getPackageInfoComments());
            }
        } catch (ParseException ex) {
            Log.e(TAG, "File skipped, error : ".concat(ex.getMessage()));
        } catch (FileNotFoundException ex) {
            Log.e(TAG, ex.getMessage());
        }
    }

    private static Language defineJavaLang() {
        Language language = new Language();
        language.setLanguage(JAVA);
        List<String> paradigms = new ArrayList<>();
        paradigms.add(OOP);
        paradigms.add(STRUCTURED);
        paradigms.add(IMPERATIVE);
        paradigms.add(GENERIC);
        paradigms.add(REFLECTIVE);
        paradigms.add(CONCURRENT);
        language.setParadigms(paradigms);

        return language;
    }

    private static String extractPackagePath(String sourceFilePath) {
        int lastSlashIndex = sourceFilePath.lastIndexOf("/");
        return sourceFilePath.substring(0, lastSlashIndex);
    }

}
