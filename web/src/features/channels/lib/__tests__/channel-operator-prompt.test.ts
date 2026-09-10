import { describe, expect, test } from 'vitest'

import {
  CHANNEL_FORM_DEFAULT_VALUES,
  buildSettingJSON,
  transformChannelToFormDefaults,
} from '../channel-form'

describe('operator system prompt channel settings', () => {
  test('buildSettingJSON persists the operator system prompt', () => {
    const formData = {
      ...CHANNEL_FORM_DEFAULT_VALUES,
      operator_system_prompt: 'operator instructions',
    }

    const setting = JSON.parse(buildSettingJSON(formData))

    expect(setting.operator_system_prompt).toBe('operator instructions')
  })

  test('transformChannelToFormDefaults reads operator_system_prompt from channel settings', () => {
    const channel = {
      channel_info: { multi_key_mode: 'random' },
      setting: '{"operator_system_prompt":"operator instructions"}',
    } as unknown as Parameters<typeof transformChannelToFormDefaults>[0]

    const defaults = transformChannelToFormDefaults(channel)

    expect(defaults.operator_system_prompt).toBe('operator instructions')
  })

  test('transformChannelToFormDefaults defaults operator_system_prompt to an empty string', () => {
    const channel = {
      channel_info: { multi_key_mode: 'random' },
      setting: '{"system_prompt":"stock"}',
    } as unknown as Parameters<
      typeof transformChannelToFormDefaults
    >[0]

    const defaults = transformChannelToFormDefaults(channel)

    expect(defaults.operator_system_prompt).toBe('')
  })
})
