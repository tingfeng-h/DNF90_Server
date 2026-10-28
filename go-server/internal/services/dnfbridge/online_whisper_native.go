package dnfbridge

import (
	"os"
	"strings"

	dnfrepo "longheng.io/server/internal/modules/dnf/repository"
)

// 移植自「轻装上阵」分支的私聊实现：补上普通私聊（type 1 / 0x29）的完整路径。
//
// 上游只处理 0x24/0x28，普通私聊会直接回失败——这与「主输入框私聊不通」的现场一致。
// 分支的实现要点：
//   - type 1/0x29 用 mode=1（不带会话对端 DSTR），0x24/0x28 用 mode=0x24/0x28；
//   - 普通私聊的**发出方**回显用 mode=7（原生 inbound=1 / outbound=7 两种呈现）；
//   - 频道字节取会话 owner channel，而不是账号区服组。
//
//	DNFBRIDGE_WHISPER_NATIVE=on   # 默认 off，行为与上游一致
const whisperNativeEnv = "DNFBRIDGE_WHISPER_NATIVE"

func whisperNativeEnabled() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(whisperNativeEnv))) {
	case "1", "on", "true", "yes":
		return true
	default:
		return false
	}
}

func isCurrentPrivateWhisperTypeNative(messageType byte) bool {
	return messageType == 1 || messageType == 0x24 || messageType == 0x28 || messageType == 0x29
}

func buildCurrentWhisperBodyNative(mode, channel byte, sender, peer string, message []byte) []byte {
	var writer packetWriter
	writer.writeByte(mode)
	writer.writeByte(0)
	writer.writeRawDstr(rosterNameBytes(sender))
	writer.writeByte(channel)
	writer.writeRawDstr(message)
	writer.writeRawDstr(nil)
	if mode == 0x24 {
		writer.writeRawDstr(rosterNameBytes(peer))
	}
	return writer.bytes()
}

func (s *Service) sendCurrentOfficialWhisperNative(sourceSession, targetSession *gameSession, opcode uint16,
	request currentPeerChatRequest, source, target dnfrepo.CharacterRecord) error {
	if !isCurrentPrivateWhisperTypeNative(request.MessageType) {
		return s.sendGameUpperFailure(sourceSession, opcode, 1)
	}
	sourceIdentity, _ := s.boundGameSessionCharacterSnapshot(sourceSession)
	targetIdentity, _ := s.boundGameSessionCharacterSnapshot(targetSession)
	if err := s.sendGameUpperSuccess(sourceSession, opcode, nil); err != nil {
		return err
	}
	channel, _ := currentSessionOwnerChannelField(sourceSession)
	mode := byte(0x24)
	if request.MessageType == 1 || request.MessageType == 0x29 {
		mode = 1
	}
	incoming := buildCurrentWhisperBodyNative(mode, channel, source.Name, source.Name, request.Message)
	if err := s.sendGameUpperRawClassCodec(targetSession, currentOfficialWhisperNotifyOpcode, incoming, 0, true); err != nil {
		return err
	}
	s.rememberCharacterFriendChatLink(targetIdentity, source, channel)
	echoMode, echoChannel, echoName := mode, channel, source.Name
	if mode == 1 {
		// 原生普通私聊：inbound=1、outbound=7，发出那一行写的是收件人而不是自己。
		echoMode = 7
		echoName = target.Name
		echoChannel, _ = currentSessionOwnerChannelField(targetSession)
	}
	echo := buildCurrentWhisperBodyNative(echoMode, echoChannel, echoName, target.Name, request.Message)
	if err := s.sendGameUpperRawClassCodec(sourceSession, currentOfficialWhisperNotifyOpcode, echo, 0, true); err != nil {
		return err
	}
	if mode == 1 {
		s.rememberCharacterFriendChatLink(sourceIdentity, target, echoChannel)
	} else {
		s.rememberCharacterFriendChatLink(sourceIdentity, source, echoChannel)
	}
	s.logGameEvent(sourceSession, "game-official-whisper-forwarded",
		"source_char_id", selectedCharacterID(sourceSession),
		"target_char_id", selectedCharacterID(targetSession),
		"message_type", request.MessageType, "message_len", len(request.Message),
		"msg_id", currentOfficialWhisperNotifyOpcode,
		"receive_mode", mode, "echo_mode", echoMode, "echo_sent", true,
		"body_source", "fork_native_whisper_mode1_mode7_no_trade_state")
	return nil
}
