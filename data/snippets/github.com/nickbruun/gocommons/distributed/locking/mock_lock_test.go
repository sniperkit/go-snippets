package locking

import (
	logrus "github.com/Sirupsen/logrus"
	log "github.com/nickbruun/gocommons/logging"
	"github.com/nickbruun/gocommons/unittest"
	"sync"
	"testing"
)

type MockLockTestSuite struct {
	LockTestSuite

	MockLockProvider *MockLockProvider
}

func (s *MockLockTestSuite) SetUp() {
	s.MockLockProvider = NewMockLockProvider()
	s.LockProvider = s.MockLockProvider
}

func (s *MockLockTestSuite) TestFailRecover() {
	{
		l1 := s.LockProvider.GetLock("/lock/path")
		l2 := s.LockProvider.GetLock("/lock/path")

		var l1Waiting, l1Failed, l2Acquired sync.WaitGroup
		l1Waiting.Add(1)
		l1Failed.Add(1)
		l2Acquired.Add(1)

		failCh := s.AssertLock(l1)

		go func() {
			l1Waiting.Done()
			<-failCh
			l1Failed.Done()
		}()

		go func() {
			s.AssertLock(l2)
			l2Acquired.Done()
		}()

		l1Waiting.Wait()

		s.MockLockProvider.Fail()

		l1Failed.Wait()

		s.AssertUnlockNotHeld(l1)

		s.MockLockProvider.Recover()

		l2Acquired.Wait()

		l2.Unlock()
	}

	{
		l1 := s.LockProvider.GetLock("/lock/path")
		l2 := s.LockProvider.GetLock("/lock/path")

		var l1Waiting, l1Failed, l2Waiting, l2Acquired sync.WaitGroup
		l1Waiting.Add(1)
		l1Failed.Add(1)
		l2Waiting.Add(1)
		l2Acquired.Add(1)

		failCh := s.AssertLock(l1)

		go func() {
			l1Waiting.Done()
			<-failCh
			l1Failed.Done()
		}()

		l1Waiting.Wait()
		log.Debug("Lock 1 waiting to fail")

		s.MockLockProvider.Fail()
		log.Debug("Provider failed")

		l1Failed.Wait()
		log.Debug("Lock 1 failed")

		s.AssertUnlockNotHeld(l1)

		go func() {
			l2Waiting.Done()
			s.AssertLock(l2)
			l2Acquired.Done()
		}()

		l2Waiting.Wait()
		log.Debug("Lock 2 waiting to acquire")

		s.MockLockProvider.Recover()
		log.Debug("Provider recovered")

		l2Acquired.Wait()
		log.Debug("Lock 2 acquired")

		l2.Unlock()
	}
}

func TestMockLock(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	unittest.RunTestSuite(&MockLockTestSuite{}, t)
}
