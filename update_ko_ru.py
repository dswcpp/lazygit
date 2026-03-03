#!/usr/bin/env python3
import json
import os

translations_dir = r"pkg\i18n\translations"

ko_trans = {
    "AIDiffSkippedFilesNote": " (추가로 %d개의 잠금/바이너리/생성 파일 건너뜀)",
    "AIAssistantSystemPrompt": "당신은 git 명령 생성기입니다. 사용자 요구사항과 저장소 상태를 기반으로 실행해야 할 shell/git 명령을 생성합니다.\n\n",
}

for lang in ["ko"]:
    filepath = os.path.join(translations_dir, f"{lang}.json")
    with open(filepath, 'r', encoding='utf-8') as f:
        data = json.load(f)
    for key, value in ko_trans.items():
        if key in data:
            data[key] = value
    with open(filepath, 'w', encoding='utf-8') as f:
        json.dump(data, f, ensure_ascii=False, indent=2)
        f.write('\n')
    print(f"Updated {lang}.json")
