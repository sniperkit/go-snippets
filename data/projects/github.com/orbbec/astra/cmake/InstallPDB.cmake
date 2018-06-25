#
# In MSVC, installs pdb files to specified directory for Debug and RelWithDebInfo configs
#
# Usage:
#   install_pdb(YOURTARGET path/in/installdir/
#

function(install_pdb _targetname _destdir) 
    if(NOT MSVC)
        return()
    endif()
    
    if(NOT TARGET ${_targetname})
        message("install_pdb: There is no target named '${_targetname}'")
        return()
    endif()
    
    install(FILES "$<TARGET_FILE_DIR:${_targetname}>/${_targetname}${CMAKE_DEBUG_POSTFIX}.pdb"
            DESTINATION "${_destdir}"
            CONFIGURATIONS Debug
            )
    install(FILES "$<TARGET_FILE_DIR:${_targetname}>/${_targetname}.pdb"
            DESTINATION "${_destdir}"
            CONFIGURATIONS RelWithDebInfo
            )
endfunction()
