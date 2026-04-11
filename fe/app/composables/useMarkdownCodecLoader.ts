type MarkdownCodecModule = typeof import('~/utils/markdownCodec')

let markdownCodecPromise: Promise<MarkdownCodecModule> | null = null

export const useMarkdownCodecLoader = () => {
  const ensureMarkdownCodec = async () => {
    if (!markdownCodecPromise) {
      markdownCodecPromise = import('~/utils/markdownCodec')
    }
    return markdownCodecPromise
  }

  return {
    ensureMarkdownCodec
  }
}
