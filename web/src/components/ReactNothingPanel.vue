<template>
  <div ref="host" class="react-host"></div>
</template>

<script setup>
import { onMounted, onBeforeUnmount, ref, watch } from 'vue'
import React from 'react'
import { createRoot } from 'react-dom/client'
import { NothingCard, DotMatrixText, ProgressDots } from '@dennislee928/nothingx-react-components'

const props = defineProps({
  title: { type: String, default: 'XDEVICE LAB' },
  subtitle: { type: String, default: 'Cross-protocol exchange confidence' },
  value: { type: Number, default: 82 }
})

const host = ref(null)
let root

const render = () => {
  if (!root) return
  root.render(
    React.createElement(
      NothingCard,
      {
        dark: true,
        style: {
          width: '100%',
          background: '#101010',
          color: '#ffffff',
          borderRadius: '20px',
          border: '1px solid #2d2d2d',
          boxShadow: '0 16px 40px rgba(0, 0, 0, 0.35)'
        }
      },
      React.createElement(DotMatrixText, { children: props.title, color: '#ffffff', dotSize: 3, gap: 1 }),
      React.createElement('div', { style: { marginTop: 10, fontFamily: 'IBM Plex Mono, monospace', fontSize: 12, color: '#9ca3af' } }, props.subtitle),
      React.createElement(
        'div',
        { style: { marginTop: 12 } },
        React.createElement(ProgressDots, { value: Math.max(0, Math.min(100, props.value)), max: 100, count: 24, color: '#ff4d4f' })
      )
    )
  )
}

onMounted(() => {
  root = createRoot(host.value)
  render()
})

watch(() => [props.title, props.subtitle, props.value], render)

onBeforeUnmount(() => {
  if (root) root.unmount()
})
</script>

<style scoped>
.react-host {
  width: 100%;
}
</style>
