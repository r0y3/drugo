package main

import (
	"errors"
	"testing"
)

func TestNew(t *testing.T) {
	svc := New()

	if svc == nil {
		t.Error("New should not return nil")
	} else {
		if svc.registry == nil {
			t.Fatal("Drupal service registry should not be nil")
		}
		if svc.nodeService == nil {
			t.Fatal("Drupal service nodeService should not be nil")
		}
		if svc.fetched == nil {
			t.Fatal("Drupal service fetched channel should not be nil")
		}
		if svc.done == nil {
			t.Fatal("Drupal service done channel should not be nil")
		} else {
			select {
			case svc.done <- true:
				break
			default:
				t.Fatal("Should be able to send data to done channel.")
			}
		}
		if svc.err == nil {
			t.Fatal("Drupal service err channel should not be nil")
		} else {
			select {
			case svc.err <- nil:
				break
			default:
				t.Fatal("Should be able to send data to err channel.")
			}
		}
	}
}

func TestFetch(t *testing.T) {
	t.Log("It should be implemented.")
}

func TestSave(t *testing.T) {
	t.Log("It should be implemented.")
}

func TestDone(t *testing.T) {
	svc := New()

	svc.done <- true
	select {
	case done := <-svc.Done():
		if !done {
			t.Fatal("It should be done.")
		}
	}
	svc.done <- false
	select {
	case done := <-svc.Done():
		if done {
			t.Fatal("It should not be done.")
		}
	}
}

func TestError(t *testing.T) {
	svc := New()

	svc.err <- nil
	select {
	case err := <-svc.Error():
		if err != nil {
			t.Fatal("It should be nil.")
		}
	}
	testErr := errors.New("Test error.")
	svc.err <- testErr
	select {
	case err := <-svc.Error():
		if err != testErr {
			t.Fatalf("It should be '%s'.", testErr)
		}
	}
}
