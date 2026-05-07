<template>
  <q-layout view="hHh lpR fFf">
    <q-header class="bg-white text-black" bordered>
      <q-toolbar>
        <q-toolbar-title class="text-weight-bold">XDevice InfoCard Mission Console</q-toolbar-title>
        <q-chip color="red-5" text-color="white">Quasar + React Components</q-chip>
      </q-toolbar>
    </q-header>

    <q-page-container>
      <q-page class="q-pa-md q-pa-lg-xl" style="max-width: 1200px; margin: 0 auto;">
        <q-card flat bordered class="q-mb-lg" style="background: var(--paper)">
          <q-card-section>
            <div class="text-h5 text-weight-bold">What these 4 projects try to achieve</div>
            <div class="text-body1 text-grey-8 q-mt-sm">
              A portfolio of contact-exchange approaches across acoustic signaling, nearby P2P connectivity,
              visual codes, and cryptographic/on-chain attestations.
            </div>
          </q-card-section>
          <q-separator/>
          <q-card-section>
            <ReactNothingPanel />
          </q-card-section>
        </q-card>

        <div class="demo-grid">
          <q-card bordered flat>
            <q-card-section>
              <div class="section-title text-red-7">01 / ULTRASONIC</div>
              <div class="text-subtitle1 text-weight-medium q-mt-xs">Audio-tone card transfer</div>
              <div class="text-body2 q-mt-sm text-grey-8">Uses shared Rust DSP to modulate/decode payload through near-ultrasonic tones.</div>
            </q-card-section>
            <q-separator />
            <q-card-section>
              <q-slider v-model="noise" label :min="0" :max="100" color="red" />
              <div class="text-caption">Demo: Estimated decode confidence {{ 100 - noise }}%</div>
            </q-card-section>
          </q-card>

          <q-card bordered flat>
            <q-card-section>
              <div class="section-title text-red-7">02 / NEARBY</div>
              <div class="text-subtitle1 text-weight-medium q-mt-xs">Device-to-device handshake</div>
              <div class="text-body2 q-mt-sm text-grey-8">Combines app state-machine + signaling to coordinate offer/answer and transfer lifecycle.</div>
            </q-card-section>
            <q-separator />
            <q-card-section>
              <q-stepper v-model="nearbyStep" flat animated>
                <q-step :name="1" title="Discover" />
                <q-step :name="2" title="Connect" />
                <q-step :name="3" title="Exchange" />
              </q-stepper>
              <q-btn dense color="red" label="Next" @click="nearbyStep = nearbyStep === 3 ? 1 : nearbyStep + 1" />
            </q-card-section>
          </q-card>

          <q-card bordered flat>
            <q-card-section>
              <div class="section-title text-red-7">03 / VISUAL HANDSHAKE</div>
              <div class="text-subtitle1 text-weight-medium q-mt-xs">Animated QR + repair symbols</div>
              <div class="text-body2 q-mt-sm text-grey-8">Frame stream plus repair chunks aims to recover payload despite dropped frames.</div>
            </q-card-section>
            <q-separator />
            <q-card-section>
              <q-input v-model.number="dropRate" type="number" label="Drop rate %" min="0" max="90" dense/>
              <q-linear-progress :value="Math.max(0, 1 - dropRate / 100)" color="red" class="q-mt-sm"/>
              <div class="text-caption">Demo recovery indicator</div>
            </q-card-section>
          </q-card>

          <q-card bordered flat>
            <q-card-section>
              <div class="section-title text-red-7">04 / WEB3 SBT</div>
              <div class="text-subtitle1 text-weight-medium q-mt-xs">Signed relationship attestations</div>
              <div class="text-body2 q-mt-sm text-grey-8">Relayer prepares EIP-712-like signable payload and tracks tx-like lifecycle metadata.</div>
            </q-card-section>
            <q-separator />
            <q-card-section>
              <q-timeline color="red" layout="dense">
                <q-timeline-entry title="Prepare" subtitle="requestId + signable payload" />
                <q-timeline-entry title="Submit" subtitle="signature checks + digest" />
                <q-timeline-entry title="Status" subtitle="track operation and tx hash" />
              </q-timeline>
            </q-card-section>
          </q-card>
        </div>
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
