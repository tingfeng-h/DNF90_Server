package dnfbridge

import (
	"errors"
	"fmt"
	"path"
	"strings"
	"sync"

	"longheng.io/server/internal/modules/dnf/actionobject"
	dnfpvf "longheng.io/server/internal/modules/dnf/pvf"
)

// 独立实现 90CN 的被动对象动作脚本解析（不修改 catalog 核心）：
// passiveobject.lst 给出 `.obj`，`.obj` 的 `[basic action]` 指向同目录 `.act`。
// 这条链驱动 [DUNGEON SYMBOL CALC] / [GATE ENABLE]，即房间机关。
const currentDungeonPassiveObjectList = "passiveobject/passiveobject.lst"

type passiveObjectScriptCache struct {
	mu      sync.Mutex
	indexed bool
	paths   map[int64]string
	loaded  map[int64]bool
	scripts map[int64]actionobject.Script
	err     error
}

func (s *Service) passiveObjectCache() *passiveObjectScriptCache {
	s.passiveObjectMu.Lock()
	defer s.passiveObjectMu.Unlock()
	if s.passiveObjectScripts == nil {
		s.passiveObjectScripts = &passiveObjectScriptCache{}
	}
	return s.passiveObjectScripts
}

func cleanDungeonPassiveObjectPath(value string) string {
	value = strings.TrimSpace(strings.ReplaceAll(value, "\\", "/"))
	value = strings.TrimPrefix(value, "./")
	return strings.TrimPrefix(value, "/")
}

func isDungeonPassiveObjectDefinitionMissing(err error) bool {
	return errors.Is(err, dnfpvf.ErrDocNotFound)
}

func readDungeonPassiveObjectText(source dnfpvf.Source, listedPath string) (string, string, error) {
	candidates := []string{cleanDungeonPassiveObjectPath(listedPath)}
	if !strings.HasPrefix(strings.ToLower(candidates[0]), "passiveobject/") {
		candidates = append(candidates, path.Join("passiveobject", candidates[0]))
	}
	var lastErr error
	for _, candidate := range candidates {
		text, err := source.ReadText(candidate)
		if err == nil {
			return candidate, text, nil
		}
		lastErr = err
	}
	return "", "", lastErr
}

func resolveDungeonActionObjectScriptPath(objectPath, basicAction string) string {
	action := strings.TrimPrefix(strings.TrimSpace(strings.ReplaceAll(basicAction, "\\", "/")), "/")
	if action == "" {
		return ""
	}
	directory := path.Dir(cleanDungeonPassiveObjectPath(objectPath))
	if directory == "." || directory == "" {
		return path.Clean(action)
	}
	if strings.HasPrefix(strings.ToLower(action), strings.ToLower(directory)+"/") {
		return path.Clean(action)
	}
	return path.Join(directory, action)
}

func (c *passiveObjectScriptCache) ensureIndexed(source dnfpvf.Source) {
	if c.indexed {
		return
	}
	c.indexed = true
	text, err := source.ReadText(currentDungeonPassiveObjectList)
	if err != nil {
		c.err = fmt.Errorf("read %s: %w", currentDungeonPassiveObjectList, err)
		return
	}
	document, err := dnfpvf.Parse(currentDungeonPassiveObjectList, text)
	if err != nil {
		c.err = fmt.Errorf("parse %s: %w", currentDungeonPassiveObjectList, err)
		return
	}
	entries := dnfpvf.ParseList(document)
	if len(entries) == 0 {
		c.err = fmt.Errorf("empty %s", currentDungeonPassiveObjectList)
		return
	}
	c.paths = make(map[int64]string, len(entries))
	for _, entry := range entries {
		if _, exists := c.paths[entry.ID]; !exists {
			c.paths[entry.ID] = cleanDungeonPassiveObjectPath(entry.Path)
		}
	}
}

func (s *Service) passiveObjectActionScript(source dnfpvf.Source, id int64) (actionobject.Script, bool, error) {
	if source == nil {
		return actionobject.Script{}, false, errors.New("passive object source unavailable")
	}
	c := s.passiveObjectCache()
	c.mu.Lock()
	defer c.mu.Unlock()
	c.ensureIndexed(source)
	if c.err != nil {
		return actionobject.Script{}, false, c.err
	}
	if c.loaded == nil {
		c.loaded = make(map[int64]bool)
	}
	if c.loaded[id] {
		script, found := c.scripts[id]
		return script, found, nil
	}
	c.loaded[id] = true
	listedPath, listed := c.paths[id]
	if !listed {
		return actionobject.Script{}, false, nil
	}
	objectPath, objectText, err := readDungeonPassiveObjectText(source, listedPath)
	if err != nil {
		if isDungeonPassiveObjectDefinitionMissing(err) {
			return actionobject.Script{}, false, nil
		}
		return actionobject.Script{}, false, fmt.Errorf("read passive object id=%d path=%s: %w", id, listedPath, err)
	}
	objectDocument, err := dnfpvf.Parse(objectPath, objectText)
	if err != nil {
		return actionobject.Script{}, false, fmt.Errorf("parse passive object id=%d path=%s: %w", id, objectPath, err)
	}
	basicAction, _ := objectDocument.Text("basic action")
	scriptPath := resolveDungeonActionObjectScriptPath(objectPath, basicAction)
	if scriptPath == "" {
		return actionobject.Script{}, false, nil
	}
	scriptText, err := source.ReadText(scriptPath)
	if err != nil {
		if isDungeonPassiveObjectDefinitionMissing(err) {
			return actionobject.Script{}, false, nil
		}
		return actionobject.Script{}, false, fmt.Errorf("read action script id=%d path=%s: %w", id, scriptPath, err)
	}
	scriptDocument, err := dnfpvf.Parse(scriptPath, scriptText)
	if err != nil {
		return actionobject.Script{}, false, fmt.Errorf("parse action script id=%d path=%s: %w", id, scriptPath, err)
	}
	script, err := actionobject.Parse(scriptDocument)
	if err != nil {
		return actionobject.Script{}, false, fmt.Errorf("decode action script id=%d path=%s: %w", id, scriptPath, err)
	}
	if !script.Drives() {
		return actionobject.Script{}, false, nil
	}
	if c.scripts == nil {
		c.scripts = make(map[int64]actionobject.Script)
	}
	c.scripts[id] = script
	return script, true, nil
}
