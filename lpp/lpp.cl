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

(defstruct param type name deref funcset)
(defstruct funcdef funcset funcname params returntype)

(setq preprocessor-file-extension "lpp")
(setq begin-marker-left "^^^BEGINREPLACE:")
(setq begin-marker-right "^^^")
(setq end-marker "^^^ENDREPLACE^^^")
(setq auto-header-marker "^^^AUTOGENHEADER^^^")
(setq token-marker #\^)
(setq token-arg-marker #\:)
(setq token-arg-delimiter #\,)

(defstruct lppmacro macro filter)

(setq macro-hash (make-hash-table :test 'equal))
(defun add-macro (&rest args)
  (let ((m (apply #'make-lppmacro args)))
    (setf (gethash (lppmacro-macro m) macro-hash) m)))
(defun get-macro (key) (gethash key macro-hash))

(setq void-param-hash (make-hash-table :test 'equal))
(defun add-void-param (&rest args)
  (let* ((p (apply #'make-param args))
         (funcset-name (param-funcset p)))
    (cond ((not (null funcset-name))
           (setf (gethash funcset-name void-param-hash) p))
          (t (error "Must provide :funcset for add-void-param")))))

(defun get-void-param (funcset-name) (gethash funcset-name void-param-hash))

(add-macro :macro "RETURN"      :filter (lambda (fd len args vp) (funcdef-returntype fd)))
(add-macro :macro "FUNC"        :filter (lambda (fd len args vp) (funcdef-funcname fd)))
(add-macro :macro "PARAMS"      :filter (lambda (fd len args vp) (format-params (funcdef-params fd) len args vp)))

;; PARAMS arguments: types, names, deref, ref, void, voidonly, nowrap
;;;;;;;;;;;;;;;;;;

(defun partial (func &rest args1)
  (lambda (&rest args2) (apply func (append args1 args2))))

(defun newline ()
  (write-string (format nil "~%")))

(defun concat-string-list (s)
  (format nil "~{~A~^~%~}" s))

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

(defun back-to-forward-slashes (str)
  (replace-all str "\\" "/")
)

(defun is-arg-set (target-arg args)
  (not (null (find target-arg args :test 'equal))))

(defun format-paramitem-type (args p)
  (format nil "~A" (param-type p)))

(defun format-paramitem-name (args p)
  ;;in name only format, we can dereference the argument (pointer)
  ;;but only if the macro argument requested it
  ;;and this parameter can be dereferenced
  (if (and (is-arg-set "deref" args) (param-deref p))
      ;;then
      (format nil "*~A" (param-name p))
      ;;else
      (format nil "~A" (param-name p))))

(defun format-paramitem-full (args p)
  (let* ((type-decl (param-type p))
         (last-index (1- (length type-decl))))
    (format nil "~A ~A"
            ;;if macro requested reference mode,
            ;; this parameter is dereferencable,
            ;; and the last char is a *
            (if (and (is-arg-set "ref" args) (param-deref p) (equal "*" (subseq type-decl last-index)))
                ;;then, replace the last * with a &
                (concatenate 'string (subseq type-decl 0 last-index) "&")
                ;;else
                type-decl)
            (param-name p))))

(defun error-if-null (var errmsg) (when (null var) (error errmsg)))

(defun format-params (params token-start-length args void-param)
  ;;full, types, names, deref, ref, void, voidonly, nowrap
  (let*  ((arg-types      (is-arg-set "types" args))
          (arg-names      (is-arg-set "names" args))
          (arg-deref      (is-arg-set "deref" args))
          (arg-ref        (is-arg-set "ref" args))
          (arg-void       (is-arg-set "void" args))
          (arg-voidonly   (is-arg-set "voidonly" args))
          (arg-wrap       (is-arg-set "wrap" args))
          (arg-nowrap     (is-arg-set "nowrap" args))
          (filtered-params (cond  (arg-voidonly
                                   (error-if-null void-param "Void param must be set to use void argument")
                                   (cons void-param nil)) ;;only the void-param
                                  (arg-void
                                   (error-if-null void-param "Void param must be set to use void argument")
                                   (cons void-param params)) ;;prepend the void-param
                                  (t params))) ;;default, just use params
          (types-only (and arg-types (not arg-names)))
          (names-only (and arg-names (not arg-types)))
          (full-decl  (not (xor arg-names arg-types)))
          (format-param-func (cond  (full-decl  #'format-paramitem-full)
                                    (types-only #'format-paramitem-type)
                                    (names-only #'format-paramitem-name)))
          (do-wrap    (or arg-wrap (not (or names-only arg-nowrap))))
          (format-line (if do-wrap "~{~A~^,~%~}" "~{~A~^, ~}")))
    (when filtered-params
      (format nil format-line
              (cons
               ;;first parameter => no indentation
               (apply format-param-func (list args (car filtered-params)))
               (mapcar
                (lambda (p)
                  (concatenate 'string
                               ;;second parameter, indent according to wrap policy
                               (when do-wrap (make-string token-start-length :initial-element #\ ))
                               (apply format-param-func (list args p))))
                (cdr filtered-params)))))))

(defun split-by-delimiter (string delimiter)
  (loop   for start = 0 then (1+ finish)
     for finish = (position delimiter string :start start)
     if (and (not (eq (length string) start))
             (not (eq start finish)))
     collect (subseq string start finish)
     until (null finish)))

(defun process-tokens (line funcdata void-param &optional (line-length 0))
  (let ((token-marker-left (position token-marker line)))
    (if (null token-marker-left)
        ;;then, no left token marker. return line
        line
        ;;else, process more
        (let* ((token-start (1+ token-marker-left))
               (token-stop (position token-marker line :start token-start)))
          (if (null token-stop)
              ;;then, no right token marker. return line
              line
              ;;else, process more
              (let* ((token-full (subseq line token-start token-stop))
                     (seperator-index (position token-arg-marker token-full))
                     (token (subseq token-full 0 seperator-index))
                     (token-args (if (not (null seperator-index))
                                     (split-by-delimiter
                                      (subseq token-full (1+ seperator-index) nil) token-arg-delimiter)
                                     '(nil)))
                     (line-begin (subseq line 0 token-marker-left))
                     (line-end (subseq line (1+ token-stop) nil))
                     (m (get-macro token)))
                ;;m (the macro struct) is nil if no macro found for the token
                (if (null m)
                    ;;then
                    line
                    ;;else
                    (let* ((token-start-length (+ line-length token-marker-left))
                           (expanded-token
                            (funcall (lppmacro-filter m)
                                     funcdata
                                     token-start-length
                                     token-args void-param)))
                      (concatenate 'string
                                   line-begin
                                   expanded-token
                                   (process-tokens
                                    line-end
                                    funcdata
                                    void-param
                                    (+ token-start-length (length expanded-token))))))))))))

(defun expand-template (template-list func-data void-param)
  (concat-string-list
   (mapcar
    (lambda (template-line) (process-tokens template-line func-data void-param))
    template-list)))

(defun expand-methods-with-template (template-list funcset-name)
  (remove nil ;;filter out the nil values (from wrong funcsets)
          (mapcar
           (lambda (func-data)
             (cond ;;only expand this function if it belongs to the target funcset
               ((equal funcset-name (funcdef-funcset func-data))
                (expand-template template-list func-data (get-void-param funcset-name)))
               ;;otherwise return nil
               (t nil)))
           funcs)))

(defun clear-funcs ()
  (setq funcs nil)
  (setq funcs-tail funcs))

(clear-funcs)

(defun add-func (&rest args)
  (let ((func-data (apply #'make-funcdef args)))
    (if (endp funcs)
        ;;then
        (progn
          (setq funcs (cons func-data nil))
          (setq funcs-tail funcs))
        ;;else
        (progn
          (setf (cdr funcs-tail) (cons func-data nil))
          (setq funcs-tail (cdr funcs-tail))))))

(defun outputfuncs (f)
  (write-line (concat-string-list f)))

(defun filter-line (line filename)
  (replace-all line
               auto-header-marker
               (format nil "THIS FILE AUTO-GENERATED FROM ~A. DO NOT EDIT." filename)))

(defmacro prepend-into (var value)
  `(setq ,var (cons ,value ,var)))

(defun strip-end-characters (str numChars)
  (let ((index (- (length str) numChars)))
    (if (>= index 1)
        (subseq str 0 index)
        nil)))

(defun get-funcset-from-begin-marker (line)
  (let ((start-index (search begin-marker-left line)))
    (cond ;;no begin marker, return nil
      ((null start-index) nil)

      ;;begin marker must be at start of line
      ((not (= start-index 0)) nil)

      ;;left side is good
      (t (let*  ((start-param-index (+ start-index (length begin-marker-left)))
                 (end-index (search begin-marker-right line :start2 start-param-index)))
           (cond ;;no right side marker, return nil
             ((null end-index) nil)

             ;;right marker must be at end of line
             ((not (= (length line) (+ end-index (length begin-marker-right)))) nil)

             ;;markers but no funcset parameter, can't process
             ((= start-param-index end-index) nil)
             ;;we have a begin-marker and funcset parameter!
             (t (values (subseq line start-param-index end-index) start-param-index end-index))))))))

(defun is-end-marker (line)
  (equal line end-marker))

(defun in-to-out-filename (fn)
  (strip-end-characters fn 4))

(defun process-file (infilename)
  (let ((infile (open infilename :if-does-not-exist nil))
        (outfile (open (in-to-out-filename infilename) :direction :output :if-exists :supersede))
        (filelines nil)
        (template-lines nil)
        (funcset-name nil))
    (when infile
      (loop
         for line = (read-line infile nil)
         while line
         do
           (cond   ;;outside of a replace block
             ((null funcset-name)
              ;;look for a begin-replace marker (funcset name becomes not nil)
              (setq funcset-name (get-funcset-from-begin-marker line))
              ;;if still no funcset-name
              (when (null funcset-name)
                ;;accumulate filtered lines to filelines
                (prepend-into filelines (filter-line line (file-namestring infilename)))))

             ;;all below: funcset-name is not nil -- inside of a replace block
             ;;check for end marker
             ((is-end-marker line)
              ;;since we have a full template, expand template it
              (mapcar
               (lambda (v) (prepend-into filelines v))
               (expand-methods-with-template (reverse template-lines) funcset-name))
              (setq funcset-name nil)  ;;clear funcset-name -- return to normal operation above
              (setq template-lines nil) ;;clear template so we can process another one in the same file
             )
             ;;inside replace block, accumulate template
             (t (prepend-into template-lines line))))
      (loop for line in (reverse filelines)
         do (write-line line outfile))
      (close infile)
      (close outfile))))

(defvar api-file "apidef.cl")
(defvar lpp-file "lpp.cl")

(load api-file :verbose nil)
(load "cl-fad/load.lisp" :verbose nil)

(defun mapc-directory-tree (fn directory)
  (dolist (entry (cl-fad:list-directory directory))
    (when (cl-fad:directory-pathname-p entry)
      (mapc-directory-tree fn entry))
    (funcall fn entry)))

(defun _ () (load lpp-file :verbose nil))

(defun file-modified-time (f)
  (cond ((probe-file f) (posix:file-stat-mtime (posix:file-stat f)))
        (t 0)
  )
)

(defun process (dir)
  (let* ((target-directory (if (null dir)
                               (ext:cd)
                               dir))
         (cache-file-path (cl-fad:merge-pathnames-as-file
                           (cl-fad:pathname-as-directory target-directory)
                           ".pp-modification-cache.cl"))
         (api-modify (file-modified-time api-file))
         (lpp-modify (file-modified-time lpp-file))
         (pp-cache (with-open-file (fsi cache-file-path :if-does-not-exist nil)
                      (cond ((null fsi) (make-hash-table :test 'equal) )
                            (t (let ((data (read fsi nil)))
                                  (if (hash-table-p data)
                                      data
                                      (make-hash-table :test 'equal)
                                   )
                               ))
                      )
                   ))
        )
    (write-line (format nil "=== LPP processing directory: ~A ===~%" target-directory))
    (mapc-directory-tree (lambda (x)
                           (when (equal (pathname-type x) preprocessor-file-extension)
                             (let* ((fwd-path (back-to-forward-slashes (namestring x)))
                                    (cache-modify (gethash fwd-path pp-cache 0))
                                    (last-modify (max lpp-modify api-modify (file-modified-time x)))
                                   )
                               (when (> last-modify cache-modify)
                                 (write-line (format nil "~A => ~A"
                                                     (enough-namestring x target-directory)
                                                     (in-to-out-filename (file-namestring x))
                                             ))

                                 (process-file (namestring x))
                                 (setf (gethash fwd-path pp-cache) last-modify)
                               )
                             )
                           )
                         )
                         target-directory)
    (with-open-file (fso cache-file-path :direction :output :if-does-not-exist :create :if-exists :supersede)
      (write pp-cache :stream fso)
    )
  )
)

(process (car *args*))
;(process "../include")
;(process "../src")
