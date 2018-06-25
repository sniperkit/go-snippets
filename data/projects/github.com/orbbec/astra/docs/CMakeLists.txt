find_package(Sphinx REQUIRED)

if(NOT DEFINED SPHINX_THEME)
  set(SPHINX_THEME custom_sphinx_rtd_theme)
endif()

if (NOT DEFINED SPHINX_THEME_DIR)
  set(SPHINX_THEME_DIR "${CMAKE_CURRENT_SOURCE_DIR}/_themes")
endif()

set(ASTRA_DOCS_BUILDERS "html" CACHE STRING "Build Astra Documentation")

set(BINARY_BUILD_DIR "${CMAKE_CURRENT_BINARY_DIR}/_build")

set(SPHINX_CACHE_DIR "${CMAKE_CURRENT_BINARY_DIR}/_doctrees")

set(SPHINX_HTML_DIR "${CMAKE_CURRENT_BINARY_DIR}/html")

configure_file(
  "${CMAKE_CURRENT_SOURCE_DIR}/conf.py.in"
  "${BINARY_BUILD_DIR}/conf.py"
  @ONLY)

set(BUILD_HTML_DOCS false)

foreach(_builder ${ASTRA_DOCS_BUILDERS})
   if(${_builder} STREQUAL "html")
      set(BUILD_HTML_DOCS true)
   else()
      message(WARNING "${_builder} is not a supported documentation builder type. Ignoring.")
   endif()
endforeach(_builder)

set(SPHINX_HTML_COMMAND 
   ${SPHINX_EXECUTABLE}
    -q -b html
    -c "${BINARY_BUILD_DIR}"
    -d "${SPHINX_CACHE_DIR}"
    "${CMAKE_CURRENT_SOURCE_DIR}"
    "${SPHINX_HTML_DIR}"
)

if(${BUILD_HTML_DOCS})
  add_custom_target(AstraSDKDocumentation ALL
    COMMAND ${SPHINX_HTML_COMMAND}
    COMMENT "Building HTML documentation with Sphinx")
endif()

set_target_properties(AstraSDKDocumentation PROPERTIES FOLDER "documentation")

install(DIRECTORY ${SPHINX_HTML_DIR}
  DESTINATION docs
  PATTERN ".buildinfo" EXCLUDE
  PATTERN "objects.inv" EXCLUDE
  )