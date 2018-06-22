/*
 * Copyright 2014-2015 The DevMine Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package ch.devmine.javaparser.utils;

import java.text.DateFormat;
import java.text.SimpleDateFormat;
import java.util.Date;

/**
 *
 * @author gilliek
 */
public class Log {

    public static final String DATE_TIME_FORMAT = "yyyy/MM/dd HH:mm:ss ";

    /**
     * The output format of the log.
     */
    private static final String LOG_FORMAT = "%s[%s]\t%s\t%s\n";

    /**
     * Current log level.
     */
    private static int logLevel = 6;

    /**
     * Priority constant for verbose logging; use Log.v.
     */
    public static final int VERBOSE = 2;

    /**
     * Priority constant for debug logging; use Log.d.
     */
    public static final int DEBUG = 3;

    /**
     * Priority constant for verbose logging; use Log.i.
     */
    public static final int INFO = 4;

    /**
     * Priority constant for warning logging; use Log.w.
     */
    public static final int WARN = 5;

    /**
     * Priority constant for error logging; use Log.e.
     */
    public static final int ERROR = 6;

    /**
     * Send a VERBOSE log message.
     *
     * @param tag Used to identify the source of a log message. It usually
     * identifies the class or activity where the log call occurs.
     * @param message The message you would like logged.
     */
    public static void v(String tag, String message) {
        println(VERBOSE, tag, message);
    }

    /**
     * Send a VERBOSE log message.
     *
     * @param tag Used to identify the source of a log message. It usually
     * identifies the class or activity where the log call occurs.
     * @param message The message you would like logged.
     * @param tr An exception to log.
     */
    public static void v(String tag, String message, Throwable tr) {
        println(VERBOSE, tag, message + "\n" + getStackTraceString(tr));
    }

    /**
     * Send a DEBUG log message.
     *
     * @param tag Used to identify the source of a log message. It usually
     * identifies the class or activity where the log call occurs.
     * @param message The message you would like logged.
     */
    public static void d(String tag, String message) {
        println(DEBUG, tag, message);
    }

    /**
     * Send a DEBUG log message.
     *
     * @param tag Used to identify the source of a log message. It usually
     * identifies the class or activity where the log call occurs.
     * @param message The message you would like logged.
     * @param tr An exception to log.
     */
    public static void d(String tag, String message, Throwable tr) {
        println(DEBUG, tag, message + "\n" + getStackTraceString(tr));
    }

    /**
     * Send a INFO log message.
     *
     * @param tag Used to identify the source of a log message. It usually
     * identifies the class or activity where the log call occurs.
     * @param message The message you would like logged.
     */
    public static void i(String tag, String message) {
        println(INFO, tag, message);
    }

    /**
     * Send a INFO log message.
     *
     * @param tag Used to identify the source of a log message. It usually
     * identifies the class or activity where the log call occurs.
     * @param message The message you would like logged.
     * @param tr An exception to log.
     */
    public static void i(String tag, String message, Throwable tr) {
        println(INFO, tag, message + "\n" + getStackTraceString(tr));
    }

    /**
     * Send a WARN log message.
     *
     * @param tag Used to identify the source of a log message. It usually
     * identifies the class or activity where the log call occurs.
     * @param message The message you would like logged.
     */
    public static void w(String tag, String message) {
        println(WARN, tag, message);
    }

    /**
     * Send a WARN log message.
     *
     * @param tag Used to identify the source of a log message. It usually
     * identifies the class or activity where the log call occurs.
     * @param message The message you would like logged.
     * @param tr An exception to log.
     */
    public static void w(String tag, String message, Throwable tr) {
        println(WARN, tag, message + "\n" + getStackTraceString(tr));
    }

    /**
     * Send a ERROR log message.
     *
     * @param tag Used to identify the source of a log message. It usually
     * identifies the class or activity where the log call occurs.
     * @param message The message you would like logged.
     */
    public static void e(String tag, String message) {
        println(ERROR, tag, message);
    }

    /**
     * Send a ERROR log message.
     *
     * @param tag Used to identify the source of a log message. It usually
     * identifies the class or activity where the log call occurs.
     * @param message The message you would like logged.
     * @param tr An exception to log.
     */
    public static void e(String tag, String message, Throwable tr) {
        println(ERROR, tag, message + "\n" + getStackTraceString(tr));
    }

    /**
     * Handy function to get a loggable stack trace from a Throwable
     *
     * @param tr An exception to log
     * @return The stack trace as a "loggable" String.
     */
    public static String getStackTraceString(Throwable tr) {
        if (tr == null) {
            return "";
        }

        StringBuilder builder = new StringBuilder();

        for (StackTraceElement ste : tr.getStackTrace()) {
            builder.append("\t");
            builder.append(ste.getClassName());
            builder.append(":");
            builder.append(ste.getLineNumber());
            builder.append("\t");
            builder.append(ste.getMethodName());
            builder.append("\n");
        }

        return builder.toString();
    }

    /**
     * Getter for the log level. The default value is Log.ERROR-
     *
     * @return The current log level.
     */
    public static int getLogLevel() {
        return logLevel;
    }

    /**
     * Set the log level.
     *
     * @param level New log level
     */
    public static void setLogLevel(int level) {
        logLevel = level;
    }

    /**
     * Print the log according to its level.
     *
     * @param level Level of the log message.
     * @param tag Used to identify the source of a log message. It usually
     * identifies the class or activity where the log call occurs.
     * @param message The message you would like logged.
     */
    private static void println(int level, String tag, String message) {
       // Only log if the supplied level is less than or equal to the current
        // log level.
        if (level > logLevel) {
            return;
        }

        DateFormat dateFormat = new SimpleDateFormat(DATE_TIME_FORMAT);
        Date now = new Date();
        String fullMessage = String.format(LOG_FORMAT, dateFormat.format(now),
            levelToString(level), tag, message);

        switch (level) {
            case VERBOSE:
            case DEBUG:
            case INFO:
            case WARN:
                System.err.print(fullMessage);
                break;
            case ERROR:
                System.err.print(fullMessage);
                break;
           default:
           }
    }

    /**
     * Get the level name corresponding to a level number.
     *
     * @param level The level to convert
     * @return The level name as a String.
     */
    private static String levelToString(int level) {
        switch (level) {
            case VERBOSE:
                return "VERBOSE";
            case DEBUG:
                return "DEBUG";
            case INFO:
                return "INFO";
            case WARN:
                return "WARN";
            case ERROR:
                return "ERROR";
            default:
                return "UNKNOWN";
        }
    }
}
