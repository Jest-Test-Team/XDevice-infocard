<template>
  <q-layout view="hHh lpR fFf" class="app-shell">
    <q-header class="shell-header" bordered>
      <q-toolbar class="toolbar-wrap">
        <div class="brand-wrap reveal r1">
          <div class="brand-kicker">XDEVICE INFOSYSTEM</div>
          <q-toolbar-title class="brand-title">InfoCard Transfer Lab</q-toolbar-title>
        </div>
        <q-chip dense class="build-chip reveal r2">Quasar + NothingX React</q-chip>
      </q-toolbar>
    </q-header>

    <q-page-container>
      <q-page class="page-wrap">
        <section class="hero-grid panel panel-hero reveal r2">
          <div>
            <div class="hero-eyebrow">MISSION OVERVIEW</div>
            <h1 class="hero-title">Four contact-exchange projects, one reusable trust flow.</h1>
            <p class="hero-subtitle">
              A practical portfolio comparing acoustic payloads, nearby P2P state machines,
              visual frame-recovery, and signed on-chain style attestations.
            </p>
            <div class="hero-metrics">
              <div class="metric-card">
                <div class="metric-label">Implementations</div>
                <div class="metric-value">4</div>
              </div>
              <div class="metric-card">
                <div class="metric-label">Core protocol families</div>
                <div class="metric-value">Audio / Nearby / Visual / Web3</div>
              </div>
            </div>
          </div>
          <ReactNothingPanel class="hero-react" title="SYSTEM STATE" subtitle="Cross-protocol reliability" :value="82" />
        </section>

        <section class="projects-grid">
          <q-card bordered flat class="project-card reveal r2">
            <q-card-section>
              <div class="project-kicker">01 / ULTRASONIC</div>
              <h2 class="project-title">Audio-tone card transfer</h2>
              <p class="project-copy">Uses shared Rust DSP to modulate/decode payload through near-ultrasonic tones.</p>
            </q-card-section>
            <q-separator />
            <q-card-section>
              <q-slider v-model="noise" :min="0" :max="100" color="dark" label label-always/>
              <div class="demo-row"><span>Noise injection</span><strong>{{ noise }}%</strong></div>
              <div class="demo-row"><span>Decode confidence</span><strong>{{ 100 - noise }}%</strong></div>
            </q-card-section>
          </q-card>

          <q-card bordered flat class="project-card reveal r3">
            <q-card-section>
              <div class="project-kicker">02 / NEARBY</div>
              <h2 class="project-title">Device-to-device handshake</h2>
              <p class="project-copy">Coordinates discovery, signaling, offer/answer, and payload transfer lifecycle.</p>
            </q-card-section>
            <q-separator />
            <q-card-section>
              <q-stepper v-model="nearbyStep" flat animated class="mini-stepper">
                <q-step :name="1" title="Discover" />
                <q-step :name="2" title="Connect" />
                <q-step :name="3" title="Exchange" />
              </q-stepper>
              <q-btn flat color="dark" class="step-btn" label="Advance state" @click="nearbyStep = nearbyStep === 3 ? 1 : nearbyStep + 1" />
            </q-card-section>
          </q-card>

          <q-card bordered flat class="project-card reveal r4">
            <q-card-section>
              <div class="project-kicker">03 / VISUAL HANDSHAKE</div>
              <h2 class="project-title">Animated QR + repair symbols</h2>
              <p class="project-copy">Streams frame chunks plus repair symbols to recover payload under dropped frames.</p>
            </q-card-section>
            <q-separator />
            <q-card-section>
              <q-input v-model.number="dropRate" type="number" label="Frame drop rate %" min="0" max="90" dense outlined />
              <q-linear-progress :value="Math.max(0, 1 - dropRate / 100)" color="dark" class="q-mt-md" size="10px" />
              <div class="demo-row q-mt-sm"><span>Recovery estimate</span><strong>{{ Math.max(10, 100 - dropRate) }}%</strong></div>
            </q-card-section>
          </q-card>

          <q-card bordered flat class="project-card reveal r5">
            <q-card-section>
              <div class="project-kicker">04 / WEB3 SBT</div>
              <h2 class="project-title">Signed relationship attestations</h2>
              <p class="project-copy">Prepares EIP-712-style signable payload, relays it, and tracks operation state.</p>
            </q-card-section>
            <q-separator />
            <q-card-section>
              <q-timeline color="grey-9" layout="dense">
                <q-timeline-entry title="Prepare" subtitle="requestId + signable payload" />
                <q-timeline-entry title="Submit" subtitle="signature checks + digest" />
                <q-timeline-entry title="Track" subtitle="operation id + tx hash" />
              </q-timeline>
            </q-card-section>
          </q-card>
        </section>

        <section class="secondary-react-grid">
          <div class="reveal r4"><ReactNothingPanel title="ULTRASONIC" subtitle="Live decode confidence" :value="100 - noise" /></div>
          <div class="reveal r5"><ReactNothingPanel title="VISUAL" subtitle="Repair pipeline health" :value="Math.max(10, 100 - dropRate)" /></div>
        </section>
      </q-page>
    </q-page-container>
  </q-layout>
</template>

<script setup>
import { ref } from 'vue'
import ReactNothingPanel from './components/ReactNothingPanel.vue'

const noise = ref(18)
const nearbyStep = ref(1)
const dropRate = ref(25)
</script>
