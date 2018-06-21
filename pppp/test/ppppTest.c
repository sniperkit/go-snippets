#include "pppp.h"
#include <llog.h>
#include <random.h>
#include <unity.h>
#include "fault.h"
#include "unity_internals.h"
#include "leak_detector.h"


void test_failure_func(void)
{
      TEST_ASSERT_EQUAL_STRING("ello", "Hello");
}

void test_success_func(void)
{
      TEST_ASSERT_EQUAL_STRING("Hello", "Hello");
}

void test_ignore_func(void)
{
      TEST_IGNORE_MESSAGE("This is a ignored test!");
}

int main(int argc, char** argv)
{
    register_fault_handlers();
    memleak_detect_start();
    llog(INFO, "%s", "Start the test for Project tt!");
    UNITY_BEGIN();

    RUN_TEST(test_failure_func);
    RUN_TEST(test_success_func);
    RUN_TEST(test_ignore_func);

    UNITY_END();
    llog(INFO, "%s", "Test For Project tt is Finished!");
    memleak_detect_stop();
    return 0;
 }
