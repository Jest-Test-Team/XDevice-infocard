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
        dark: false,
        style: {
          width: '100%',
          background: '#ffffff',
          color: '#0b0b0b',
          borderRadius: '24px',
          border: '1px solid rgba(0, 0, 0, 0.10)',
          boxShadow: '0 20px 60px rgba(0, 0, 0, 0.06)'
        }
      },
      React.createElement(DotMatrixText, { children: props.title, color: '#0b0b0b', dotSize: 3, gap: 1 }),
      React.createElement('div', { style: { marginTop: 10, fontFamily: 'IBM Plex Mono, monospace', fontSize: 12, color: 'rgba(0, 0, 0, 0.55)' } }, props.subtitle),
      React.createElement(
        'div',
        { style: { marginTop: 12 } },
        React.createElement(ProgressDots, { value: Math.max(0, Math.min(100, props.value)), max: 100, count: 24, color: '#ff2d2d' })
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
