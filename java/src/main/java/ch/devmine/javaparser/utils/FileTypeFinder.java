/*
 * Copyright 2014-2015 The DevMine Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package ch.devmine.javaparser.utils;

import java.io.File;
import java.io.FilenameFilter;

/**
 *
 * @author lweingart
 */
public class FileTypeFinder {


    /**
     * class tag
     */
    public static final String TAG = "FileTypeFinder::";

    /**
     *
     * @param directory
     *      the path of the searched directory
     * @param fileType
     *      the searched file extension (ex: .java)
     * @return
     *      a boolean if the directory contains the searched file type
     */
    public static boolean search(String directory, String fileType) {
        return new FileTypeFinder().listFile(directory, fileType);
    }

    /**
     *
     * @param folder
     *      the path of the searched folder
     * @param ext
     *      the searched file extension (ex: .java)
     * @return
     *      a boolean if the directory contains the searched file type
     */
    public boolean listFile(String folder, String ext) {

        GenericExtFilter filter = new GenericExtFilter(ext);
        File dir = new File(folder);
        if(!dir.isDirectory()){
            Log.i(TAG, "Directory does not exists : " + folder);
            return false;
        }

        // list out all the file name and filter by the extension
        String[] list = dir.list(filter);

        if (list.length == 0) {
//            Log.i(TAG, "folder " + folder + " doesn't contain any file ending with : " + ext);
            return false;
        } else {
            return true;
        }
    }

    /**
     * inner class, generic extension filter
     */
    public class GenericExtFilter implements FilenameFilter {

        /**
         * file extension
         */
        private final String ext;

        /**
         * construct the generic extension filter
         * @param ext
         *      the file extension
         */
        public GenericExtFilter(String ext) {
            this.ext = ext;
        }

        @Override
        public boolean accept(File dir, String name) {
            return (name.endsWith(ext));
        }
    }
}
