/*
* Copyright 2014-2015 The DevMine Authors. All rights reserved.
* Use of this source code is governed by a BSD-style
* license that can be found in the LICENSE file.
*/
package ch.devmine.javaparser.utils;

import java.util.ArrayList;
import java.util.List;

/**
 *
 * @author lweingart
 */
public class FileWalkerUtils {

    /**
     *
     * @param fileList
     *      a list of files
     * @return
     *      a list of contained java files
     */
    public static List<String> extractJavaFiles(List<String> fileList) {
        List<String> result = new ArrayList<>();
        for (String string : fileList) {
            if (string.endsWith(".java")) {
                result.add(string);
            }
        }
        return result;
    }

    /**
     *
     * @param filePath
     *      path of a directory
     * @return
     *      the name of the directory
     */
    public static String extractFolderName(String filePath) {
        int i = filePath.lastIndexOf("/");
        String folderPath = filePath.substring(0, i);
        i = folderPath.lastIndexOf("/");

        return folderPath.substring(i + 1);
    }
}
