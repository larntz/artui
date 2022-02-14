package ui

// AppsJSON is an example app description in json
var AppsJSON = `
{
    "apiVersion": "argoproj.io/v1alpha1",
    "kind": "Application",
    "metadata": {
        "annotations": {
            "kubectl.kubernetes.io/last-applied-configuration": "{\"apiVersion\":\"argoproj.io/v1alpha1\",\"kind\":\"Application\",\"metadata\":{\"annotations\":{\"notifications.argoproj.io/subscribe.slack\":\"argocd-notifications\"},\"finalizers\":[\"resources-finalizer.argocd.argoproj.io\"],\"labels\":{\"argocd.argoproj.io/instance\":\"apps\"},\"name\":\"apps\",\"namespace\":\"argocd\"},\"spec\":{\"destination\":{\"namespace\":\"argocd\",\"server\":\"https://kubernetes.default.svc\"},\"project\":\"default\",\"source\":{\"path\":\"kind/apps\",\"repoURL\":\"git@github.com:larntz/argocd-dev.git\",\"targetRevision\":\"HEAD\"},\"syncPolicy\":{\"automated\":{\"prune\":true,\"selfHeal\":true}}}}\n",
            "notifications.argoproj.io/subscribe.slack": "argocd-notifications",
            "notified.notifications.argoproj.io": "{\"4512623916d78eb8c04de558e6f5a913153cce9a:on-deployed:[0].y7b5sbwa2Q329JYH755peeq-fBs:slack:larntz-kind-testing\":1644542809}"
        },
        "creationTimestamp": "2022-02-11T01:26:06Z",
        "finalizers": [
            "resources-finalizer.argocd.argoproj.io"
        ],
        "generation": 15990,
        "labels": {
            "argocd.argoproj.io/instance": "apps"
        },
        "managedFields": [
            {
                "apiVersion": "argoproj.io/v1alpha1",
                "fieldsType": "FieldsV1",
                "fieldsV1": {
                    "f:metadata": {
                        "f:annotations": {
                            ".": {},
                            "f:notifications.argoproj.io/subscribe.slack": {}
                        },
                        "f:finalizers": {
                            ".": {},
                            "v:\"resources-finalizer.argocd.argoproj.io\"": {}
                        }
                    },
                    "f:spec": {
                        ".": {},
                        "f:destination": {
                            ".": {},
                            "f:namespace": {},
                            "f:server": {}
                        },
                        "f:project": {},
                        "f:source": {
                            ".": {},
                            "f:path": {},
                            "f:repoURL": {},
                            "f:targetRevision": {}
                        },
                        "f:syncPolicy": {
                            ".": {},
                            "f:automated": {
                                ".": {},
                                "f:prune": {},
                                "f:selfHeal": {}
                            }
                        }
                    }
                },
                "manager": "kubectl-client-side-apply",
                "operation": "Update",
                "time": "2022-02-11T01:26:06Z"
            },
            {
                "apiVersion": "argoproj.io/v1alpha1",
                "fieldsType": "FieldsV1",
                "fieldsV1": {
                    "f:metadata": {
                        "f:annotations": {
                            "f:kubectl.kubernetes.io/last-applied-configuration": {}
                        },
                        "f:labels": {
                            ".": {},
                            "f:argocd.argoproj.io/instance": {}
                        }
                    },
                    "f:status": {
                        ".": {},
                        "f:health": {
                            ".": {},
                            "f:status": {}
                        },
                        "f:history": {},
                        "f:operationState": {
                            ".": {},
                            "f:finishedAt": {},
                            "f:message": {},
                            "f:operation": {
                                ".": {},
                                "f:initiatedBy": {
                                    ".": {},
                                    "f:automated": {}
                                },
                                "f:retry": {
                                    ".": {},
                                    "f:limit": {}
                                },
                                "f:sync": {
                                    ".": {},
                                    "f:prune": {},
                                    "f:revision": {}
                                }
                            },
                            "f:phase": {},
                            "f:startedAt": {},
                            "f:syncResult": {
                                ".": {},
                                "f:resources": {},
                                "f:revision": {},
                                "f:source": {
                                    ".": {},
                                    "f:path": {},
                                    "f:repoURL": {},
                                    "f:targetRevision": {}
                                }
                            }
                        },
                        "f:reconciledAt": {},
                        "f:resources": {},
                        "f:sourceType": {},
                        "f:summary": {},
                        "f:sync": {
                            ".": {},
                            "f:comparedTo": {
                                ".": {},
                                "f:destination": {
                                    ".": {},
                                    "f:namespace": {},
                                    "f:server": {}
                                },
                                "f:source": {
                                    ".": {},
                                    "f:path": {},
                                    "f:repoURL": {},
                                    "f:targetRevision": {}
                                }
                            },
                            "f:revision": {},
                            "f:status": {}
                        }
                    }
                },
                "manager": "argocd-application-controller",
                "operation": "Update",
                "time": "2022-02-11T01:26:15Z"
            },
            {
                "apiVersion": "argoproj.io/v1alpha1",
                "fieldsType": "FieldsV1",
                "fieldsV1": {
                    "f:metadata": {
                        "f:annotations": {
                            "f:notified.notifications.argoproj.io": {}
                        }
                    }
                },
                "manager": "argocd-notifications-backend",
                "operation": "Update",
                "time": "2022-02-11T01:26:51Z"
            }
        ],
        "name": "apps",
        "namespace": "argocd",
        "resourceVersion": "1394044",
        "uid": "86ba0d08-647f-4ec4-b1a8-e4475a4e19d6"
    },
    "spec": {
        "destination": {
            "namespace": "argocd",
            "server": "https://kubernetes.default.svc"
        },
        "project": "default",
        "source": {
            "path": "kind/apps",
            "repoURL": "git@github.com:larntz/argocd-dev.git",
            "targetRevision": "HEAD"
        },
        "syncPolicy": {
            "automated": {
                "prune": true,
                "selfHeal": true
            }
        }
    },
    "status": {
        "health": {
            "status": "Healthy"
        },
        "history": [
            {
                "deployStartedAt": "2022-02-11T01:26:09Z",
                "deployedAt": "2022-02-11T01:26:15Z",
                "id": 0,
                "revision": "4512623916d78eb8c04de558e6f5a913153cce9a",
                "source": {
                    "path": "kind/apps",
                    "repoURL": "git@github.com:larntz/argocd-dev.git",
                    "targetRevision": "HEAD"
                }
            }
        ],
        "operationState": {
            "finishedAt": "2022-02-11T01:26:15Z",
            "message": "successfully synced (all tasks run)",
            "operation": {
                "initiatedBy": {
                    "automated": true
                },
                "retry": {
                    "limit": 5
                },
                "sync": {
                    "prune": true,
                    "revision": "4512623916d78eb8c04de558e6f5a913153cce9a"
                }
            },
            "phase": "Succeeded",
            "startedAt": "2022-02-11T01:26:09Z",
            "syncResult": {
                "resources": [
                    {
                        "group": "argoproj.io",
                        "hookPhase": "Running",
                        "kind": "Application",
                        "message": "application.argoproj.io/prometheus configured",
                        "name": "prometheus",
                        "namespace": "argocd",
                        "status": "Synced",
                        "syncPhase": "Sync",
                        "version": "v1alpha1"
                    },
                    {
                        "group": "argoproj.io",
                        "hookPhase": "Running",
                        "kind": "Application",
                        "message": "application.argoproj.io/apps configured",
                        "name": "apps",
                        "namespace": "argocd",
                        "status": "Synced",
                        "syncPhase": "Sync",
                        "version": "v1alpha1"
                    },
                    {
                        "group": "argoproj.io",
                        "hookPhase": "Running",
                        "kind": "Application",
                        "message": "application.argoproj.io/webapp configured",
                        "name": "webapp",
                        "namespace": "argocd",
                        "status": "Synced",
                        "syncPhase": "Sync",
                        "version": "v1alpha1"
                    },
                    {
                        "group": "argoproj.io",
                        "hookPhase": "Running",
                        "kind": "Application",
                        "message": "application.argoproj.io/argocd configured",
                        "name": "argocd",
                        "namespace": "argocd",
                        "status": "Synced",
                        "syncPhase": "Sync",
                        "version": "v1alpha1"
                    },
                    {
                        "group": "argoproj.io",
                        "hookPhase": "Running",
                        "kind": "Application",
                        "message": "application.argoproj.io/traefik configured",
                        "name": "traefik",
                        "namespace": "argocd",
                        "status": "Synced",
                        "syncPhase": "Sync",
                        "version": "v1alpha1"
                    }
                ],
                "revision": "4512623916d78eb8c04de558e6f5a913153cce9a",
                "source": {
                    "path": "kind/apps",
                    "repoURL": "git@github.com:larntz/argocd-dev.git",
                    "targetRevision": "HEAD"
                }
            }
        },
        "reconciledAt": "2022-02-14T18:17:04Z",
        "resources": [
            {
                "group": "argoproj.io",
                "kind": "Application",
                "name": "apps",
                "namespace": "argocd",
                "status": "Synced",
                "version": "v1alpha1"
            },
            {
                "group": "argoproj.io",
                "kind": "Application",
                "name": "argocd",
                "namespace": "argocd",
                "status": "Synced",
                "version": "v1alpha1"
            },
            {
                "group": "argoproj.io",
                "kind": "Application",
                "name": "prometheus",
                "namespace": "argocd",
                "status": "Synced",
                "version": "v1alpha1"
            },
            {
                "group": "argoproj.io",
                "kind": "Application",
                "name": "traefik",
                "namespace": "argocd",
                "status": "Synced",
                "version": "v1alpha1"
            },
            {
                "group": "argoproj.io",
                "kind": "Application",
                "name": "webapp",
                "namespace": "argocd",
                "status": "Synced",
                "version": "v1alpha1"
            }
        ],
        "sourceType": "Helm",
        "summary": {},
        "sync": {
            "comparedTo": {
                "destination": {
                    "namespace": "argocd",
                    "server": "https://kubernetes.default.svc"
                },
                "source": {
                    "path": "kind/apps",
                    "repoURL": "git@github.com:larntz/argocd-dev.git",
                    "targetRevision": "HEAD"
                }
            },
            "revision": "4512623916d78eb8c04de558e6f5a913153cce9a",
            "status": "Synced"
        }
    }
}
`
