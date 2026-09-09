/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

For commercial licensing, please contact support@quantumnous.com
*/
import { useEffect, useMemo, useState } from 'react'
import { useTranslation } from 'react-i18next'

import { RichContent } from '@/components/rich-content'
import { Button } from '@/components/ui/button'
import { ScrollArea } from '@/components/ui/scroll-area'
import { useAnnouncements } from '@/features/dashboard/hooks/use-status-data'
import type { AnnouncementItem } from '@/features/dashboard/types'

const STORAGE_KEY = 'fullscreen_announcement_seen'

function announcementKey(a: AnnouncementItem): string {
  return `${a.id ?? ''}:${a.content}`
}

/**
 * Shows the latest announcement as a full-screen markdown overlay, once per
 * session. Closing it stores the announcement key in sessionStorage so it does
 * not re-appear until the announcement changes or a new session starts.
 */
export function FullscreenAnnouncement() {
  const { t } = useTranslation()
  const { items } = useAnnouncements()
  const [open, setOpen] = useState(false)

  const latest = useMemo<AnnouncementItem | null>(
    () => items[0] ?? null,
    [items]
  )

  useEffect(() => {
    if (!latest) return
    try {
      const seen = sessionStorage.getItem(STORAGE_KEY)
      if (seen === announcementKey(latest)) return
      setOpen(true)
    } catch {
      setOpen(true)
    }
  }, [latest])

  if (!open || !latest) return null

  const dismiss = () => {
    setOpen(false)
    try {
      sessionStorage.setItem(STORAGE_KEY, announcementKey(latest))
    } catch {
      /* storage unavailable */
    }
  }

  return (
    <div className='fixed inset-0 z-[9999] flex items-center justify-center bg-black/80 p-4 backdrop-blur-sm sm:p-8'>
      <div className='flex max-h-full w-full max-w-3xl flex-col overflow-hidden rounded-2xl border bg-background shadow-2xl'>
        <div className='flex items-center justify-between gap-4 border-b px-6 py-4'>
          <div className='flex items-center gap-2'>
            <span className='bg-primary/10 text-primary inline-flex size-8 items-center justify-center rounded-full text-lg'>
              📢
            </span>
            <h2 className='text-lg font-semibold'>
              {t('Announcement')}
            </h2>
          </div>
          <Button variant='ghost' size='sm' onClick={dismiss}>
            {t('Close')} ✕
          </Button>
        </div>
        <ScrollArea className='min-h-0 flex-1 px-6 py-5'>
          <div className='prose prose-sm dark:prose-invert max-w-none'>
            <RichContent content={latest.content} />
          </div>
        </ScrollArea>
        <div className='flex justify-end gap-3 border-t px-6 py-4'>
          <Button onClick={dismiss}>{t('Got it')}</Button>
        </div>
      </div>
    </div>
  )
}
