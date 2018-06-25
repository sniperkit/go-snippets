if (DEFINED CMAKE_SCRIPT_MODE_FILE)

  if(NOT EXISTS ${DEST})
    execute_process(COMMAND ${CMAKE_COMMAND} -E copy ${SRC} ${DEST})
  endif()

else()

  function(copy_if_not_exist tgt src dest)

    if (NOT IS_ABSOLUTE ${src})
      set(src "${CMAKE_CURRENT_SOURCE_DIR}/${src}")
    endif()

    add_custom_command(
      TARGET ${tgt}
      POST_BUILD
      COMMAND ${CMAKE_COMMAND} -DDEST=${dest} -DSRC=${src} -P ${CMAKE_SOURCE_DIR}/cmake/CopyIfNotExist.cmake)

  endfunction()

endif()
