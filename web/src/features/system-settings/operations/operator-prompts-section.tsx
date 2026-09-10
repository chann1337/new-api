import { Plus, Trash2 } from 'lucide-react'
import { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'

import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'

import { SettingsForm } from '../components/settings-form-layout'
import { SettingsPageFormActions } from '../components/settings-page-context'
import { SettingsSection } from '../components/settings-section'
import { useUpdateOption } from '../hooks/use-update-option'

type OperatorPromptRow = {
  id: string
  model: string
  prompt: string
}

type OperatorPromptsSectionProps = {
  defaultPrompts: string
}

function createRow(model: string, prompt: string): OperatorPromptRow {
  return { id: crypto.randomUUID(), model, prompt }
}

function parsePrompts(json: string): OperatorPromptRow[] {
  let parsed: Record<string, string> = {}
  try {
    parsed = JSON.parse(json || '{}')
  } catch {
    parsed = {}
  }
  if (typeof parsed !== 'object' || parsed === null || Array.isArray(parsed)) {
    return []
  }
  return Object.entries(parsed).map(([model, prompt]) =>
    createRow(model, String(prompt))
  )
}

export function OperatorPromptsSection(props: OperatorPromptsSectionProps) {
  const { t } = useTranslation()
  const updateOption = useUpdateOption()
  const [rows, setRows] = useState<OperatorPromptRow[]>(() =>
    parsePrompts(props.defaultPrompts)
  )

  useEffect(() => {
    setRows(parsePrompts(props.defaultPrompts))
  }, [props.defaultPrompts])

  const isDirty =
    JSON.stringify(
      rows.map(({ model, prompt }) => ({ model, prompt }))
    ) !==
    JSON.stringify(
      parsePrompts(props.defaultPrompts).map(({ model, prompt }) => ({
        model,
        prompt,
      }))
    )

  const updateRow = (index: number, patch: Partial<OperatorPromptRow>) => {
    setRows((prev) =>
      prev.map((row, i) => (i === index ? { ...row, ...patch } : row))
    )
  }

  const resetRows = () => setRows(parsePrompts(props.defaultPrompts))

  const onSave = async () => {
    const prompts: Record<string, string> = {}
    rows.forEach((row) => {
      const model = row.model.trim()
      if (model) {
        prompts[model] = row.prompt
      }
    })
    await updateOption.mutateAsync({
      key: 'ModelSystemPrompts',
      value: JSON.stringify(prompts),
    })
  }

  return (
    <SettingsSection title={t('Operator System Prompts')}>
      <SettingsForm
        onSubmit={(event) => {
          event.preventDefault()
          onSave()
        }}
      >
        <SettingsPageFormActions
          onSave={onSave}
          isSaving={updateOption.isPending}
          isSaveDisabled={!isDirty}
          onReset={resetRows}
          isResetDisabled={!isDirty}
        />
        <p className='text-muted-foreground text-sm'>
          {t(
            'Operator prompts are always injected as the highest priority system instruction for the selected model, even when the client sends its own system prompt. Use * as the model name to apply a prompt to all models.'
          )}
        </p>
        <div className='flex flex-col gap-3'>
          {rows.map((row, index) => (
            <div
              key={row.id}
              className='grid items-start gap-2 md:grid-cols-[220px_1fr_auto]'
            >
              <Input
                value={row.model}
                placeholder={t('Model name or *')}
                onChange={(event) =>
                  updateRow(index, { model: event.target.value })
                }
              />
              <Textarea
                rows={2}
                value={row.prompt}
                placeholder={t('System prompt for this model')}
                onChange={(event) =>
                  updateRow(index, { prompt: event.target.value })
                }
              />
              <Button
                type='button'
                variant='ghost'
                size='icon'
                aria-label={t('Remove prompt')}
                onClick={() =>
                  setRows((prev) => prev.filter((_, i) => i !== index))
                }
              >
                <Trash2 className='h-4 w-4' />
              </Button>
            </div>
          ))}
          <div>
            <Button
              type='button'
              variant='outline'
              size='sm'
              onClick={() => setRows((prev) => [...prev, createRow('', '')])}
            >
              <Plus data-icon='inline-start' />
              <span>{t('Add model prompt')}</span>
            </Button>
          </div>
        </div>
      </SettingsForm>
    </SettingsSection>
  )
}
