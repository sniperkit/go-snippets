/*
 * Copyright 2014-2015 The DevMine Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package ch.devmine.javaparser;

import java.io.IOException;
import java.nio.file.FileVisitResult;
import java.nio.file.Path;
import java.nio.file.SimpleFileVisitor;
import java.nio.file.attribute.BasicFileAttributes;
import java.util.ArrayList;
import java.util.List;
import static java.nio.file.FileVisitResult.CONTINUE;
import static java.nio.file.FileVisitResult.SKIP_SUBTREE;

/**
 *
 * @author lweingart
 */
public class FilesWalker extends SimpleFileVisitor<Path> {

    /**
     * class tag
     */
    public static final String TAG = "FilesWalker class::";
    private final List<String> regularFiles;
    private final List<String> directories;

    /**
     * Construct the FileWalker
     */
    public FilesWalker() {
        super();
        regularFiles = new ArrayList<>();
        directories = new ArrayList<>();
    }

    /*
     * Warning, if the parser is called on the current directory, give the full
     * path as argument instead of only "." as it would be skipped with the
     * current implementation
     */
    @Override
    public FileVisitResult preVisitDirectory(Path dir, BasicFileAttributes attr) {
        if (dir.getFileName().toString().startsWith(".") ||
                dir.getFileName().toString().startsWith("~") ||
                dir.getFileName().toString().endsWith("~")) {
            return SKIP_SUBTREE;
        }
        return CONTINUE;
    }

    /*
     * Print information about each type of file.
     */
    @Override
    public FileVisitResult visitFile(Path file, BasicFileAttributes attr) {
        if (attr.isRegularFile()) {
            regularFiles.add(String.format("%s", file));
        }
        return CONTINUE;
    }

    /*
     * Print each directory visited.
     */
    @Override
    public FileVisitResult postVisitDirectory(Path dir, IOException exc) {
        directories.add(String.format("%s%n", dir));
        return CONTINUE;
    }

    /*
     * If there is some error accessing the file, the user is informed.
     * If this method is not overriden and an error occurs, an IOException
     * is thrown
     */
    @Override
    public FileVisitResult visitFileFailed(Path file, IOException exc) {
        System.err.println(exc);
        return CONTINUE;
    }

    /**
     *
     * @return regular files
     */
    public List<String> getRegularFiles() {
        return this.regularFiles;
    }

    /**
     *
     * @return directories
     */
    public List<String> getDirectories() {
        return this.directories;
    }
}
