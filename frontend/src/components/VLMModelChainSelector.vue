<template>
  <div class="vlm-chain-selector">
    <div
      v-for="(modelId, index) in localModelIds"
      :key="`vlm-chain-${index}`"
      class="vlm-chain-row"
    >
      <div class="vlm-chain-row__label">
        <span class="vlm-chain-row__badge">{{ index + 1 }}</span>
        <span>{{ index === 0 ? primaryLabel : fallbackLabel(index) }}</span>
        <span v-if="index === 0" class="required">*</span>
      </div>
      <div class="vlm-chain-row__selector">
        <ModelSelector
          model-type="VLLM"
          :selected-model-id="modelId"
          :all-models="allModels"
          :status="index === 0 ? status : 'default'"
          :placeholder="index === 0 ? primaryPlaceholder : fallbackPlaceholder"
          @update:selected-model-id="(value: string) => updateModel(index, value)"
          @add-model="handleAddModel"
        />
      </div>
      <div class="vlm-chain-row__actions">
        <t-button
          size="small"
          variant="text"
          :disabled="index === 0"
          @click="moveModel(index, -1)"
        >
          上移
        </t-button>
        <t-button
          size="small"
          variant="text"
          :disabled="index === localModelIds.length - 1"
          @click="moveModel(index, 1)"
        >
          下移
        </t-button>
        <t-button
          v-if="index > 0"
          size="small"
          variant="text"
          theme="danger"
          @click="removeModel(index)"
        >
          删除
        </t-button>
      </div>
    </div>

    <div class="vlm-chain-footer">
      <t-button
        size="small"
        variant="outline"
        :disabled="localModelIds.length >= maxModels"
        @click="addFallback"
      >
        添加备用模型
      </t-button>
      <span class="vlm-chain-footer__hint">
        按顺序尝试，最多 {{ maxModels }} 个；报错、超时或空白响应会自动切到下一个。
      </span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { MessagePlugin } from 'tdesign-vue-next'
import ModelSelector from '@/components/ModelSelector.vue'

const props = withDefaults(defineProps<{
  modelIds?: string[]
  allModels?: any[]
  maxModels?: number
  status?: 'default' | 'success' | 'warning' | 'error'
  primaryLabel?: string
  fallbackLabelPrefix?: string
  primaryPlaceholder?: string
  fallbackPlaceholder?: string
}>(), {
  modelIds: () => [''],
  allModels: () => [],
  maxModels: 5,
  status: 'default',
  primaryLabel: '主 VLLM 模型',
  fallbackLabelPrefix: '备用模型',
  primaryPlaceholder: '选择主 VLLM 模型',
  fallbackPlaceholder: '选择备用 VLLM 模型',
})

const emit = defineEmits<{
  'update:modelIds': [value: string[]]
  'add-model': []
}>()

const localModelIds = computed(() => {
  const ids = Array.isArray(props.modelIds) && props.modelIds.length > 0
    ? props.modelIds
    : ['']
  return ids.slice(0, props.maxModels)
})

function fallbackLabel(index: number) {
  return `${props.fallbackLabelPrefix} ${index}`
}

function normalize(ids: string[]) {
  const result: string[] = []
  const seen = new Set<string>()
  for (const raw of ids) {
    const id = (raw || '').trim()
    if (id && seen.has(id)) continue
    if (id) seen.add(id)
    result.push(id)
    if (result.length >= props.maxModels) break
  }
  return result.length > 0 ? result : ['']
}

function updateModel(index: number, value: string) {
  const next = [...localModelIds.value]
  const trimmed = (value || '').trim()
  if (trimmed && next.some((id, i) => i !== index && id === trimmed)) {
    MessagePlugin.warning('该 VLLM 模型已在顺序列表中。')
    return
  }
  next[index] = trimmed
  emit('update:modelIds', normalize(next))
}

function addFallback() {
  if (localModelIds.value.length >= props.maxModels) return
  emit('update:modelIds', [...localModelIds.value, ''])
}

function handleAddModel() {
  emit('add-model')
}

function removeModel(index: number) {
  if (index <= 0) return
  const next = localModelIds.value.filter((_, i) => i !== index)
  emit('update:modelIds', normalize(next))
}

function moveModel(index: number, delta: number) {
  const nextIndex = index + delta
  if (nextIndex < 0 || nextIndex >= localModelIds.value.length) return
  const next = [...localModelIds.value]
  const current = next[index]
  next[index] = next[nextIndex]
  next[nextIndex] = current
  emit('update:modelIds', normalize(next))
}
</script>

<style scoped lang="less">
.vlm-chain-selector {
  display: flex;
  flex-direction: column;
  gap: 12px;
  width: 100%;
}

.vlm-chain-row {
  display: grid;
  grid-template-columns: 120px minmax(0, 1fr) auto;
  gap: 12px;
  align-items: center;
}

.vlm-chain-row__label {
  display: flex;
  align-items: center;
  gap: 6px;
  color: var(--td-text-color-primary);
  font-size: 13px;
  white-space: nowrap;
}

.vlm-chain-row__badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 22px;
  height: 22px;
  border-radius: 999px;
  background: var(--td-brand-color-light);
  color: var(--td-brand-color);
  font-weight: 600;
}

.vlm-chain-row__selector {
  min-width: 0;
}

.vlm-chain-row__actions {
  display: flex;
  align-items: center;
  gap: 4px;
}

.vlm-chain-footer {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.vlm-chain-footer__hint {
  color: var(--td-text-color-secondary);
  font-size: 12px;
}

.required {
  color: var(--td-error-color);
}

@media (max-width: 768px) {
  .vlm-chain-row {
    grid-template-columns: 1fr;
  }

  .vlm-chain-row__actions {
    justify-content: flex-start;
  }
}
</style>
