;; This file is part of the Orbbec Astra SDK [https://orbbec3d.com]
;; Copyright (c) 2015 Orbbec 3D
;;
;; Licensed under the Apache License, Version 2.0 (the "License");
;; you may not use this file except in compliance with the License.
;; You may obtain a copy of the License at
;;
;; http://www.apache.org/licenses/LICENSE-2.0
;;
;; Unless required by applicable law or agreed to in writing, software
;; distributed under the License is distributed on an "AS IS" BASIS,
;; WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
;; See the License for the specific language governing permissions and
;; limitations under the License.
;;
;; Be excellent to each other.

(use-package :regexp)

(load "cl-fad/load.lisp" :verbose nil)

(defstruct rdata start-marker end-marker replacement-callback include-markers)

(setq project-file-extension "vcxproj")
(setq solution-file-extension "sln")
(setq filters-file-extension "filters")

(defun regexp-replace (string pat repl)
  (reduce #'(lambda (x y) (string-concat x repl y))
          (regexp:regexp-split pat string))
)

(defun file-string (path)
  (with-open-file (infile path :if-does-not-exist nil)
    (cond ((not (null infile))
            (format nil "~{~A~^~%~}"
              (loop
                for line = (read-line infile nil)
                while line
                collect line
              )
            )
          )
          (t nil) ; file not found
    )
  )
)

(defun write-file (file-name content)
  (with-open-file (outfile file-name
      :direction :output
      :if-exists :supersede
      :if-does-not-exist :create)
  (format outfile "~A" content)))

(defmacro get-next-substr-index (str &key start-marker end-marker start-index end-index include-markers (next-index 0))
  `(progn
    (let ((start-m (match ,start-marker ,str :start ,next-index))
         )
      (cond ((null start-m) nil)
            (t
              (setq ,start-index  (if include-markers
                                     (match-start start-m)
                                     (match-end start-m)
                                  )
              )
              (let ((end-m (match ,end-marker ,str :start (match-end start-m))))
                (cond ((null end-m) nil)
                      (t
                        (setq ,end-index  (if include-markers
                                             (match-end end-m)
                                             (match-start end-m)
                                          )
                        )
                        (match-end end-m)
                      )
                )
              )
            )
      )
    )
  )
)
(defun concat-string-list (s)
  (let ((CUSTOM:*PPRINT-FIRST-NEWLINE* nil))
    (format nil "~{~A~}" s)
  )
)

(defun parse-replace (str template)
    (concat-string-list
      (let ((sm (rdata-start-marker template))
            (em (rdata-end-marker template))
            (callback (rdata-replacement-callback template))
            (include-markers (rdata-include-markers template))
            )
           ;(write-line (format nil "Replacing: ~A to ~A" sm em))
           (loop  with index = 0
                  do (setq next-index (get-next-substr-index str :start-marker sm
                               :end-marker em
                               :start-index si
                               :end-index ei
                               :next-index index
                               :include-markers include-markers))
                  if (null next-index)
                    collect (progn ;(write-line (format nil "~A -> ~A" index (length str)))
                              (subseq str index)
                             )
                  while next-index
                  collect (subseq str index si)
                  collect (funcall callback (subseq str si ei))
                  collect (subseq str ei next-index)
                  do (progn ;(write-line (format nil "~A -> [~A ~A] -> ~A" index si ei next-index))
                        (setq index next-index)
                      )
           )
      )
    )
)

(defun apply-replacement-templates (str templates)
  (cond ((null templates) str)
        (t (apply-replacement-templates
            (parse-replace str (car templates))
            (cdr templates)
           )
        )
  )
)

(defun apply-parse-file (file-name-in file-name-out templates)
  (write-file file-name-out
    (apply-replacement-templates (file-string file-name-in) templates)
  )
)

(defun replace-all (string part replacement &key (test #'char=))
  "Returns a new string in which all the occurences of the part
is replaced with replacement."
  (with-output-to-string (out)
    (loop with part-length = (length part)
       for old-pos = 0 then (+ pos part-length)
       for pos = (search part string
                         :start2 old-pos
                         :test test)
       do (write-string string out
                        :start old-pos
                        :end (or pos (length string)))
       when pos do (write-string replacement out)
       while pos)))

(defmacro create-rdata-list (listname addname)
  `(progn
    (setq ,listname nil)
    (defun ,addname (&rest args)
      (let ((rd (apply #'make-rdata args)))
        (setq ,listname (cons rd ,listname))
      )
    )
  )
)

(create-rdata-list proj-data-list add-projdata)
(create-rdata-list sln-data-list add-slndata)
(create-rdata-list filter-data-list add-filterdata)

(setq tp1 (make-rdata :start-marker "<data[^>]*>"
                     :end-marker "</data>"
                     :replacement-callback (lambda (str) "_123_")
                     :include-markers T))

(defun back-to-forward-slashes (str)
  (replace-all str "\\" "/")
)

(setq lambda-build-path-to-sln-dir (lambda (str) (regexp-replace (back-to-forward-slashes str)
                                                                 (back-to-forward-slashes build-dir)
                                                                 "$(SolutionDir)")))
(setq lambda-source-path-to-source-dir (lambda (str) (regexp-replace (back-to-forward-slashes str)
                                                                     (back-to-forward-slashes source-dir)
                                                                     "$(SolutionDir)..")))
(setq lambda-convert-include-path (lambda (str) (funcall lambda-build-path-to-sln-dir
                                                         (funcall lambda-source-path-to-source-dir str)
                                                         )))
(setq lambda-int-dir-to-obj-dir (lambda (str) (regexp-replace (back-to-forward-slashes str)
                                                              "[^.]*[.]dir"
                                                              "obj")))
(setq lambda-remove (lambda (str) ""))

(add-projdata :start-marker "<AdditionalIncludeDirectories>"
              :end-marker "</AdditionalIncludeDirectories>"
              :replacement-callback lambda-convert-include-path)
(add-projdata :start-marker "<OutDir[^>]*>"
              :end-marker "</OutDir>"
              :replacement-callback lambda-build-path-to-sln-dir)
(add-projdata :start-marker "<ImportLibrary[^>]*>"
              :end-marker "</ImportLibrary>"
              :replacement-callback lambda-build-path-to-sln-dir)
(add-projdata :start-marker "<ProgramDataBaseFile[^>]*>"
              :end-marker "</ProgramDataBaseFile>"
              :replacement-callback lambda-build-path-to-sln-dir)
(add-projdata :start-marker "[^<>]*<ItemGroup[^>]*>[^<]*<CustomBuild[^>]*CMakeLists[.]txt[^>]*>"
              :end-marker "</ItemGroup>"
              :replacement-callback lambda-remove
              :include-markers T)
(add-projdata :start-marker "[^<>]*<ItemGroup[^>]*>[^<]*<ProjectReference[^>]*ZERO_CHECK[.]vcxproj[^>]*>"
              :end-marker "</ItemGroup>"
              :replacement-callback lambda-remove
              :include-markers T)
(add-projdata :start-marker "CMAKE_INTDIR="
              :end-marker ";"
              :replacement-callback lambda-remove
              :include-markers T)
(add-projdata :start-marker "<Cl\\(Include\\|Compile\\)[^iI>]*Include=\""
              :end-marker "\""
              :replacement-callback lambda-source-path-to-source-dir)
(add-projdata :start-marker "<IntDir[^>]*>"
              :end-marker "<"
              :replacement-callback lambda-int-dir-to-obj-dir)

(add-filterdata :start-marker "<Cl\\(Include\\|Compile\\)[^iI>]*Include=\""
              :end-marker "\""
              :replacement-callback lambda-source-path-to-source-dir)
(add-filterdata :start-marker "[^<>]*<ItemGroup[^>]*>[^<]*<CustomBuild[^>]*CMakeLists[.]txt[^>]*>"
              :end-marker "</ItemGroup>"
              :replacement-callback lambda-remove
              :include-markers T)

(add-slndata  :start-marker "Project([^,]*\\(ALL_BUILD\\|CMakePredefinedTargets\\|ZERO_CHECK\\|INSTALL\\)"
              :end-marker "EndProject"
              :replacement-callback lambda-remove
              :include-markers T)

(defun mapc-directory-tree (fn directory)
  (dolist (entry (cl-fad:list-directory directory))
    (when (cl-fad:directory-pathname-p entry)
      (mapc-directory-tree fn entry))
    (funcall fn entry)))

(defun get-temp-filename (fn)
  (format nil "~A.tmp" fn))

(defun process-file (filename-in file-extension)
  (let ((filename-out (get-temp-filename filename-in))
        (data-list (cond ((equal file-extension solution-file-extension) sln-data-list)
                         ((equal file-extension filters-file-extension) filter-data-list)
                         (t proj-data-list)))
       )
    (apply-parse-file filename-in filename-out data-list)
    (delete-file filename-in)
    (rename-file filename-out filename-in)
  )
)

(defun process (build-path source-path)
  (setq build-dir build-path)
  (setq source-dir source-path)
  (let ((target-directory (if (null build-dir)
                              (ext:cd)
                              build-dir)))
    (write-line (format nil "Base directory: ~A~%" target-directory))
    (mapc-directory-tree (lambda (x)
                           (let ((file-extension (pathname-type x)))
                             (when (or (equal file-extension project-file-extension)
                                       (equal file-extension filters-file-extension)
                                       (equal file-extension solution-file-extension))
                               (write-line (format nil "~A"
                                                   (enough-namestring x target-directory)
                                                   ))
                               (process-file (namestring x) file-extension))
                           )
                         )
                         target-directory)))

(defun _ () (load "parse-replace.cl" :verbose nil))

(when (>= (length *args*) 2)
  (let ((build-path (car *args*))
        (source-path (cadr *args*))
        )
    (process build-path source-path)
  )
)
