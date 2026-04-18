export type MobileHeaderAction =
  | {
      kind: "traits-copy-uuid"
      uuid: string
    }
  | {
      kind: "traits-paste-uuid"
      mode: "set" | "trait"
    }
