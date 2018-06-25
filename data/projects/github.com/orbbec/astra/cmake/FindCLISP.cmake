set(CLISP_PATHS
  "c:/Program Files (x86)/clisp-2.49/"
  "c:/Program Files/clisp-2.49/"
  /usr/local/bin)

find_program(CLISP_EXECUTABLE
  clisp
  PATHS ${CLISP_PATHS}
  DOC "Common Lisp interpreter")

include(FindPackageHandleStandardArgs)

find_package_handle_standard_args(clisp DEFAULT_MSG CLISP_EXECUTABLE)

mark_as_advanced(CLISP_EXECUTABLE)
