<template>
  <div class="kb-model-config">
    <div class="section-header">
      <h2>{{ $t('knowledgeEditor.models.title') }}</h2>
      <p class="section-description">{{ $t('knowledgeEditor.models.description') }}</p>
    </div>

    <div class="settings-group">
      <!-- LLM 大语言模型 -->
      <div class="setting-row" data-guide="kb-create-llm">
        <div class="setting-info">
          <label>{{ $t('knowledgeEditor.models.llmLabel') }} <span class="required">*</span></label>
          <p class="desc">{{ $t('knowledgeEditor.models.llmDesc') }}</p>
        </div>
        <div class="setting-control">
          <ModelSelector
            ref="llmSelectorRef"
            model-type="KnowledgeQA"
            :selected-model-id="config.llmModelId"
            :all-models="allModels"
            @update:selected-model-id="handleLLMChange"
            @add-model="handleAddModel('chat')"
            :placeholder="$t('knowledgeEditor.models.llmPlaceholder')"
          />
        </div>
      </div>

      <!-- Embedding 嵌入模型: RAG 检索启用时必填; 纯 Wiki 时可选(用于目录归类相似度) -->
      <div v-if="ragEnabled !== false || wikiEnabled" class="setting-row" data-guide="kb-create-embedding">
        <div class="setting-info">
          <label>
            {{ $t('knowledgeEditor.models.embeddingLabel') }}
            <span v-if="ragEnabled" class="required">*</span>
            <span v-else-if="wikiEnabled" class="optional">{{ $t('knowledgeEditor.models.embeddingOptional') }}</span>
          </label>
          <p class="desc">
            {{ (wikiEnabled && ragEnabled === false)
              ? $t('knowledgeEditor.models.embeddingWikiOptionalDesc')
              : $t('knowledgeEditor.models.embeddingDesc') }}
          </p>
          <t-alert
            v-if="ragEnabled && hasFiles"
            theme="warning"
            :message="$t('knowledgeEditor.models.embeddingLocked')"
            style="margin-top: 8px;"
          />
        </div>
        <div class="setting-control">
          <ModelSelector
            ref="embeddingSelectorRef"
            model-type="Embedding"
            :selected-model-id="config.embeddingModelId"
            :all-models="allModels"
            :disabled="ragEnabled && hasFiles"
            :clearable="ragEnabled === false && wikiEnabled"
            @update:selected-model-id="handleEmbeddingChange"
            @add-model="handleAddModel('embedding')"
            :placeholder="$t('knowledgeEditor.models.embeddingPlaceholder')"
          />
        </div>
      </div>

      <!-- Wiki 合成模型 (仅当 Wiki 启用时显示) -->
      <div v-if="wikiEnabled" class="setting-row">
        <div class="setting-info">
          <label>{{ $t('knowledgeEditor.wiki.synthesisModelLabel') }}</label>
          <p class="desc">{{ $t('knowledgeEditor.wiki.synthesisModelTip') }}</p>
        </div>
        <div class="setting-control">
          <VLMModelChainSelector
            :model-ids="config.wikiSynthesisModelIds || [config.wikiSynthesisModelId || '']"
            model-type="KnowledgeQA"
            :all-models="allModels"
            :primary-required="false"
            :primary-label="$t('knowledgeEditor.wiki.synthesisPrimaryLabel')"
            :fallback-label-prefix="$t('knowledgeEditor.wiki.synthesisFallbackLabelPrefix')"
            :primary-placeholder="$t('knowledgeEditor.wiki.synthesisModelPlaceholder')"
            :fallback-placeholder="$t('knowledgeEditor.wiki.synthesisFallbackPlaceholder')"
            :add-fallback-label="$t('knowledgeEditor.wiki.synthesisAddFallback')"
            :hint="$t('knowledgeEditor.wiki.synthesisModelChainHint', { max: 5 })"
            :duplicate-warning="$t('knowledgeEditor.wiki.synthesisDuplicateWarning')"
            @update:model-ids="handleWikiModelChainChange"
            @add-model="handleAddModel('knowledgeqa')"
          />
        </div>
      </div>

    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useUIStore } from '@/stores/ui'
import ModelSelector from '@/components/ModelSelector.vue'
import VLMModelChainSelector from '@/components/VLMModelChainSelector.vue'
import { useI18n } from 'vue-i18n'

interface ModelConfig {
  llmModelId?: string
  embeddingModelId?: string
  vllmModelId?: string
  wikiSynthesisModelId?: string
  wikiSynthesisModelIds?: string[]
}

interface Props {
  config: ModelConfig
  hasFiles: boolean
  wikiEnabled?: boolean
  ragEnabled?: boolean
  allModels?: any[]
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:config': [value: ModelConfig]
}>()

const uiStore = useUIStore()
const { t } = useI18n()

const llmSelectorRef = ref<InstanceType<typeof ModelSelector>>()
const embeddingSelectorRef = ref<InstanceType<typeof ModelSelector>>()

const handleLLMChange = (modelId: string) => {
  emit('update:config', {
    ...props.config,
    llmModelId: modelId
  })
}

const handleEmbeddingChange = (modelId: string) => {
  emit('update:config', {
    ...props.config,
    embeddingModelId: modelId
  })
}

const handleWikiModelChainChange = (modelIds: string[]) => {
  const normalized = normalizeModelChainIdsForEdit(modelIds)
  emit('update:config', {
    ...props.config,
    wikiSynthesisModelId: normalized[0] || '',
    wikiSynthesisModelIds: normalized
  })
}

const handleAddModel = (subSection: string) => {
  uiStore.openSettings('models', subSection)
}

const normalizeModelChainIdsForEdit = (modelIds: any[]): string[] => {
  const seen = new Set<string>()
  const ids: string[] = []

  for (const raw of modelIds || []) {
    const id = String(raw || '').trim()
    if (id) {
      if (seen.has(id)) {
        ids.push('')
        continue
      }
      seen.add(id)
      ids.push(id)
    } else {
      ids.push('')
    }
    if (ids.length >= 5) break
  }

  return ids.length > 0 ? ids : ['']
}
</script>

<style lang="less" scoped>
.kb-model-config {
  width: 100%;
}

.section-header {
  margin-bottom: 20px;

  h2 {
    font-size: 20px;
    font-weight: 600;
    color: var(--td-text-color-primary);
    margin: 0 0 6px 0;
  }

  .section-description {
    font-size: 14px;
    color: var(--td-text-color-secondary);
    margin: 0;
    line-height: 1.5;
  }
}

.settings-group {
  display: flex;
  flex-direction: column;
  gap: 0;
}

.setting-row {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  padding: 16px 0;
  border-bottom: 1px solid var(--td-component-stroke);

  &:last-child {
    border-bottom: none;
  }
}

.setting-info {
  flex: 0 0 40%;
  max-width: 40%;
  padding-right: 24px;

  label {
    font-size: 15px;
    font-weight: 500;
    color: var(--td-text-color-primary);
    display: block;
    margin-bottom: 4px;

    .required {
      color: var(--td-error-color);
      margin-left: 2px;
    }

    .optional {
      color: var(--td-text-color-placeholder);
      font-size: 12px;
      font-weight: 400;
      margin-left: 4px;
    }
  }

  .desc {
    font-size: 13px;
    color: var(--td-text-color-secondary);
    margin: 0;
    line-height: 1.5;
  }
}

.setting-control {
  flex: 0 0 55%;
  max-width: 55%;
  display: flex;
  justify-content: flex-end;
  align-items: flex-start;
}
</style>

