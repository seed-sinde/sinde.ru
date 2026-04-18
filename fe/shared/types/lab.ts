export type LabTestDepthKey =
  | "figurative"
  | "process"
  | "technology"
  | "engineering"
  | "industrial"
  | "economic"
  | "strategic"
export type LabTestRoleKey =
  | "engineer"
  | "technologist"
  | "operator"
  | "analyst"
  | "procurement"
  | "researcher"
  | "architect"
export type LabTestStageStatus = "stable" | "warning" | "critical" | "improving"
export type LabTestOption<T extends string> = {
  key: T
  label: string
}
export type LabTestTimelineStage = {
  id: string
  title: string
  short: string
  owner: string
  status: LabTestStageStatus
  duration: string
  costImpact: string
  qualityImpact: string
  bottleneck?: boolean
  details: Record<LabTestDepthKey, string>
}
export type LabTestHypothesisStatus = "draft" | "testing" | "accepted"
export type LabTestHypothesis = {
  id: string
  title: string
  stageId: string
  author: string
  effect: string
  risks: string
  status: LabTestHypothesisStatus
}
export type LabTestMetricCard = {
  label: string
  value: string
  delta: string
}
