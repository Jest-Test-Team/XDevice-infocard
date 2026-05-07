<template>
  <div ref="host"></div>
</template>

<script setup>
import { onMounted, onBeforeUnmount, ref } from 'vue'
import React from 'react'
import { createRoot } from 'react-dom/client'
import { NothingCard, DotMatrixText, ProgressDots } from '@dennislee928/nothingx-react-components'

const host = ref(null)
let root

onMounted(() => {
  root = createRoot(host.value)
  root.render(
    React.createElement(
      NothingCard,
      { dark: true, style: { maxWidth: 420, background: '#111', color: '#fff' } },
      React.createElement(DotMatrixText, { children: 'XDEVICE LAB', color: '#fff', dotSize: 3, gap: 1 }),
      React.createElement('div', { style: { marginTop: 12, fontFamily: 'monospace', fontSize: 12, color: '#aaa' } }, 'Cross-platform exchange confidence'),
      React.createElement('div', { style: { marginTop: 10 } }, React.createElement(ProgressDots, { value: 82, max: 100, count: 18, color: '#ff3b30' }))
    )
  )
})

onBeforeUnmount(() => {
  if (root) root.unmount()
})
</script>
