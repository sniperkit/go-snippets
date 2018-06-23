package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/libgit2/git2go"
	"io/ioutil"
	"os"
	"time"
)

const CURRENT_PLAYER_CONFIG_KEY string = "current_game_player"

type Event struct {
	PlayerId string
	Message  string
}

type Storage struct {
	recvEvents       <-chan string
	sendEvents       chan string
	repository       *git.Repository
	lastCommitTreeId *git.Oid
	path             string
	playerId         string
}

func NewStorage() (*Storage, error) {
	storage := &Storage{
		sendEvents: make(chan string),
	}

	return storage, nil
}

func (s *Storage) InitEventStream(events <-chan string) chan string {
	s.recvEvents = events

	go func() {
		for event := range events {
			s.storeEvent(event)
		}
	}()

	// This is the code that runs the background update process. Leaving commented out since it wont be complete for the hackathon deadline
	//go func() {
	//	for {
	//		time.Sleep(20 * time.Second)
	//		updates, err := s.getNewUpdates()
	//		if err != nil {
	//			continue
	//		}

	//		if len(updates) > 0 {
	//			s.sendEvents <- updates
	//		}
	//	}
	//}()

	go func() {
		file, err := os.Open(s.path + "/players/" + s.playerId + "/" + s.playerId + "events")
		if err != nil {
			close(s.sendEvents)
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			s.sendEvents <- scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading event stream: " + err.Error())
		}
		close(s.sendEvents)
	}()

	return s.sendEvents
}

func (s *Storage) getNewUpdates() ([]string, error) {
	remote, err := s.repository.LookupRemote("origin")
	if err != nil {
		return nil, err
	}

	err = remote.Fetch([]string{"master"}, nil, "")
	if err != nil {
		return nil, err
	}

	//originMaster, err := s.repository.LookupReference("refs/remotes/origin/master")
	//if err != nil {
	//	return nil, err
	//}

	//mergeHead, err := s.repository.AnnotatedCommitFromRef(originMaster)
	//if err != nil {
	//	return nil, err
	//}

	//mergeHeads := make([]*git.AnnotatedCommit, 1)
	//mergeHeads[0] = mergeHead
	//err = s.repository.Merge(mergeHeads, nil, nil)
	//fmt.Println("updates merged")

	head, err := s.repository.Head()
	if err != nil {
		return nil, err
	}

	headCommit, err := s.repository.LookupCommit(head.Target())
	if err != nil {
		return nil, err
	}

	headCommitTree, err := headCommit.Tree()
	if err != nil {
		return nil, err
	}

	lastTree, err := s.repository.LookupTree(s.lastCommitTreeId)
	if err != nil {
		return nil, err
	}

	diffOptions := &git.DiffOptions{
		Pathspec: []string{s.path + "/players/" + s.playerId + "/" + s.playerId + "events"},
	}
	diff, err := s.repository.DiffTreeToTree(lastTree, headCommitTree, diffOptions)

	lines := make([]git.DiffLine, 0)
	err = diff.ForEach(func(file git.DiffDelta, progress float64) (git.DiffForEachHunkCallback, error) {
		return func(hunk git.DiffHunk) (git.DiffForEachLineCallback, error) {
			return func(line git.DiffLine) error {
				lines = append(lines, line)
				return nil
			}, nil
		}, nil
	}, git.DiffDetailLines)

	updates := make([]string, 0)
	for _, line := range lines {
		if line.Origin == git.DiffLineAddition {
			updates = append(updates, line.Content)
		}
	}

	return updates, nil
}

func (s *Storage) GetCurrentPlayer() (*HeroSheet, error) {
	playerId, err := s.getFileContents(s.path + "/.git/" + CURRENT_PLAYER_CONFIG_KEY)
	if err != nil {
		return nil, err
	}
	s.playerId = string(playerId)

	contents, err := s.getFileContents(s.path + "/players/" + s.playerId + "/" + s.playerId)
	if err != nil {
		return nil, err
	}

	heroSheet := HeroSheet{}
	err = json.Unmarshal(contents, &heroSheet)
	if err != nil {
		return nil, err
	}

	return &heroSheet, nil
}

func (s *Storage) SetCurrentPlayer(playerId string) error {
	file, err := os.Create(s.path + "/.git/" + CURRENT_PLAYER_CONFIG_KEY)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(playerId)
	if err != nil {
		return err
	}

	file.Sync()
	s.playerId = playerId
	return nil
}

func (s *Storage) OpenRepository(path string) error {
	repo, err := git.OpenRepository(path)
	if err != nil {
		return err
	}

	s.repository = repo
	s.path = path

	err = s.setLastCommitTree()
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) CloneRepository(repoUrl string, path string) error {
	checkoutOptions := &git.CheckoutOpts{
		Strategy: git.CheckoutSafeCreate,
	}
	cloneOptions := &git.CloneOptions{
		Bare:           false,
		CheckoutBranch: "master",
		CheckoutOpts:   checkoutOptions,
	}

	repo, err := git.Clone(repoUrl, path, cloneOptions)
	if err != nil {
		return err
	}

	repo.CreateRemote("origin", path)
	s.repository = repo
	s.path = path

	err = s.setLastCommitTree()
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) setLastCommitTree() error {
	head, err := s.repository.Head()
	if err != nil {
		return err
	}

	headCommit, err := s.repository.LookupCommit(head.Target())
	if err != nil {
		return err
	}

	headCommitTree, err := headCommit.Tree()
	if err != nil {
		return err
	}

	s.lastCommitTreeId = headCommitTree.Id()
	return nil
}

func (s *Storage) StorePlayer(hero HeroSheet) error {
	filepath := s.path + "/players/" + hero.Name
	err := os.MkdirAll(filepath, 0755)
	if err != nil {
		return err
	}

	filename := filepath + "/" + hero.Name
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	heroBytes, err := json.Marshal(hero)
	if err != nil {
		return err
	}

	_, err = file.WriteString(string(heroBytes) + "\n")
	if err != nil {
		return err
	}

	file.Sync()
	err = s.commitCurrentIndex("Added new player: " + hero.Name)
	if err != nil {
		return err
	}

	return s.pushLatestCommits()
}

func (s *Storage) storeEvent(event string) error {
	err := os.MkdirAll(s.path+"/players/"+s.playerId, 0755)
	if err != nil {
		return err
	}

	filename := s.path + "/players/" + s.playerId + "/" + s.playerId + "events"
	file := &os.File{}
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		file, err = os.Create(filename)
		if err != nil {
			return err
		}
	} else {
		file, err = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			return err
		}
	}
	defer file.Close()

	file.WriteString(event + "\n")
	file.Sync()

	err = s.commitCurrentIndex("Event: " + event)
	if err != nil {
		return err
	}

	return s.pushLatestCommits()
}

func (s *Storage) GetGameObject(id string) ([]byte, error) {
	contents, err := s.getFileContents(id)
	if err != nil {
		return nil, err
	}

	return contents, nil
}

func (s *Storage) getFileContents(filename string) ([]byte, error) {
	contents, err := ioutil.ReadFile(filename)
	return contents, err
}

func (s *Storage) commitCurrentIndex(message string) error {
	signature := &git.Signature{
		Name:  "shark_sandwich_engine",
		Email: "shark@sandwich.com",
		When:  time.Now(),
	}

	index, err := s.repository.Index()
	if err != nil {
		return err
	}

	err = index.AddAll([]string{}, git.IndexAddDefault, nil)
	if err != nil {
		return err
	}

	treeId, err := index.WriteTree()
	if err != nil {
		return err
	}

	err = index.Write()
	if err != nil {
		return err
	}

	tree, err := s.repository.LookupTree(treeId)
	if err != nil {
		return err
	}

	head, err := s.repository.Head()
	if err != nil {
		return err
	}

	commitTarget, err := s.repository.LookupCommit(head.Target())
	if err != nil {
		return err
	}

	_, err = s.repository.CreateCommit("HEAD", signature, signature, message, tree, commitTarget)
	if err != nil {
		return err
	}

	s.lastCommitTreeId = treeId
	return nil
}

func (s *Storage) pushLatestCommits() error {
	remote, err := s.repository.LookupRemote("origin")
	if err != nil {
		return err
	}

	return remote.Push([]string{"refs/heads/master"}, nil, nil, "")
	return nil
}
