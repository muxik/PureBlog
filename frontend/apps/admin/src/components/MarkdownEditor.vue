<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref, watch } from 'vue'
import {
  Compartment,
  EditorSelection,
  EditorState,
  Prec,
  type ChangeSpec,
  type Extension,
} from '@codemirror/state'
import { defaultKeymap, history, historyKeymap } from '@codemirror/commands'
import { defaultHighlightStyle, syntaxHighlighting } from '@codemirror/language'
import { markdown } from '@codemirror/lang-markdown'
import { EditorView, keymap, placeholder as codeMirrorPlaceholder } from '@codemirror/view'

const props = withDefaults(
  defineProps<{
    modelValue: string
    placeholder?: string
    sourceMode?: boolean
  }>(),
  {
    placeholder: '',
    sourceMode: true,
  },
)

const emit = defineEmits<{
  (event: 'update:modelValue', value: string): void
  (event: 'modSlash'): void
}>()

const host = ref<HTMLDivElement | null>(null)
const placeholderCompartment = new Compartment()
const editorTheme = EditorView.theme({
  '&': {
    caretColor: 'var(--accent)',
  },
  '.cm-content': {
    caretColor: 'var(--accent)',
  },
  '.cm-cursor, .cm-dropCursor': {
    borderLeft: '2px solid var(--accent)',
  },
  '.cm-cursorLayer': {
    zIndex: '3',
  },
})

let editorView: EditorView | null = null
let syncingFromModel = false

function createPlaceholderExtension(value: string): Extension {
  return value ? codeMirrorPlaceholder(value) : []
}

function getEditorDoc(): string {
  return editorView?.state.doc.toString() ?? props.modelValue
}

function replaceDocument(value: string): void {
  if (!editorView) {
    emit('update:modelValue', value)
    return
  }

  if (value === editorView.state.doc.toString()) {
    return
  }

  editorView.dispatch({
    changes: {
      from: 0,
      to: editorView.state.doc.length,
      insert: value,
    },
  })
}

function focus(): void {
  editorView?.focus()
}

function getValue(): string {
  return getEditorDoc()
}

function setValue(value: string): void {
  replaceDocument(value)
}

function insertMarkdown(text: string): void {
  if (!editorView) {
    emit('update:modelValue', `${props.modelValue}${text}`)
    return
  }

  editorView.dispatch(editorView.state.replaceSelection(text))
  editorView.focus()
}

function wrapSelection(before: string, after: string): void {
  if (!editorView) {
    return
  }

  const { state } = editorView
  const transaction = state.changeByRange((range) => {
    const selectedText = state.doc.sliceString(range.from, range.to)

    return {
      changes: {
        from: range.from,
        to: range.to,
        insert: `${before}${selectedText}${after}`,
      },
      range: EditorSelection.range(
        range.from + before.length,
        range.from + before.length + selectedText.length,
      ),
    }
  })

  editorView.dispatch(transaction)
  editorView.focus()
}

function prefixLine(prefix: string): void {
  if (!editorView) {
    return
  }

  const { state } = editorView
  const changes: ChangeSpec[] = []
  const prefixedLines = new Set<number>()

  for (const range of state.selection.ranges) {
    const fromLine = state.doc.lineAt(range.from)
    const toPosition =
      range.to > range.from && range.to === state.doc.lineAt(range.to).from
        ? range.to - 1
        : range.to
    const toLine = state.doc.lineAt(toPosition)

    for (let lineNumber = fromLine.number; lineNumber <= toLine.number; lineNumber += 1) {
      if (prefixedLines.has(lineNumber)) {
        continue
      }

      prefixedLines.add(lineNumber)
      changes.push({
        from: state.doc.line(lineNumber).from,
        insert: prefix,
      })
    }
  }

  if (changes.length === 0) {
    return
  }

  editorView.dispatch({ changes })
  editorView.focus()
}

onMounted(() => {
  if (!host.value) {
    return
  }

  editorView = new EditorView({
    parent: host.value,
    state: EditorState.create({
      doc: props.modelValue,
      extensions: [
        markdown(),
        syntaxHighlighting(defaultHighlightStyle),
        history(),
        EditorView.lineWrapping,
        editorTheme,
        placeholderCompartment.of(createPlaceholderExtension(props.placeholder)),
        Prec.highest(
          keymap.of([
            {
              key: 'Mod-/',
              run: () => {
                emit('modSlash')
                return true
              },
            },
          ]),
        ),
        keymap.of([...defaultKeymap, ...historyKeymap]),
        EditorView.updateListener.of((update) => {
          if (!update.docChanged || syncingFromModel) {
            return
          }

          emit('update:modelValue', update.state.doc.toString())
        }),
      ],
    }),
  })
})

watch(
  () => props.modelValue,
  (value) => {
    if (!editorView || value === editorView.state.doc.toString()) {
      return
    }

    syncingFromModel = true

    try {
      editorView.dispatch({
        changes: {
          from: 0,
          to: editorView.state.doc.length,
          insert: value,
        },
      })
    } finally {
      syncingFromModel = false
    }
  },
)

watch(
  () => props.placeholder,
  (value) => {
    editorView?.dispatch({
      effects: placeholderCompartment.reconfigure(createPlaceholderExtension(value)),
    })
  },
)

onBeforeUnmount(() => {
  editorView?.destroy()
  editorView = null
})

defineExpose({
  focus,
  getValue,
  setValue,
  insertMarkdown,
  wrapSelection,
  prefixLine,
})
</script>

<template>
  <div class="markdown-editor" :data-source-mode="props.sourceMode">
    <div ref="host" class="markdown-editor__host" />
  </div>
</template>

<style scoped>
.markdown-editor {
  display: flex;
  min-height: 0;
  height: 100%;
}

.markdown-editor__host {
  flex: 1;
  min-width: 0;
  min-height: 0;
}

.markdown-editor__host :deep(.cm-editor) {
  height: 100%;
}

.markdown-editor__host :deep(.cm-scroller) {
  font-family: inherit;
  line-height: inherit;
}

.markdown-editor__host :deep(.cm-content) {
  min-height: 100%;
}
</style>
