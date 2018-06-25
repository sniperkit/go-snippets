#
# For MSVC, Creates a .vcxproj.user file with the debug working directory set.
# NOTE: Visual Studio only reads the .user when it opens the solution and will
# write it when the solution or IDE is closed. If you need to regen this file,
# close Visual Studio first and use CMake or CMake-gui.
#
# Usage:
#   set_debug_working_dir(YOURTARGET "$(TargetDir)")
#

get_filename_component(_configuserfile_dir
	${CMAKE_CURRENT_LIST_FILE}
	PATH)
set(_configuserfile_template_dir "${_configuserfile_dir}/templates")

function(set_debug_working_dir _targetname _workingDir)
    if(NOT MSVC)
        return()
    endif()
    
    if(NOT TARGET ${_targetname})
        message("set_debug_working_dir: There is no target named '${_targetname}'")
        return()
    endif()
    
    file(READ
        "${_configuserfile_template_dir}/perconfig.vcxproj.user.in"
        _userconfiggroup_template)
    
    if(CMAKE_SIZEOF_VOID_P EQUAL 8)
        set(BITS 64)
    else()
        set(BITS 32)
    endif()
    if(BITS EQUAL 64)
        set(USERFILE_PLATFORM x64)
    else()
        set(USERFILE_PLATFORM Win${BITS})
    endif()
    set(USERFILE_WORKING_DIRECTORY "${_workingDir}")
    
    foreach( USERFILE_CONFIGNAME ${CMAKE_CONFIGURATION_TYPES} )
        #string( TOUPPER ${OUTPUTCONFIG} OUTPUTCONFIG )
        string(CONFIGURE "${_userconfiggroup_template}"
                         _userconfiggroup
                         @ONLY ESCAPE_QUOTES)
        if (USERFILE_CONFIGSECTIONS)
            set(USERFILE_CONFIGSECTIONS
                "${USERFILE_CONFIGSECTIONS}\n${_userconfiggroup}")
        else()
            set(USERFILE_CONFIGSECTIONS
                "${_userconfiggroup}")
        endif()
    endforeach( USERFILE_CONFIGNAME CMAKE_CONFIGURATION_TYPES )

    set(VCPROJNAME "${CMAKE_CURRENT_BINARY_DIR}/${_targetname}")
    configure_file("${_configuserfile_template_dir}/vcxproj.user.in"
                    ${VCPROJNAME}.vcxproj.user
                    @ONLY)

endfunction()