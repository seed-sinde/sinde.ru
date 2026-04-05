import { defineEventHandler } from 'h3'
import { renderModelViewerPage } from '../../utils/modelViewer'
export default defineEventHandler(event => {
  return renderModelViewerPage(event, '3D model')
})
