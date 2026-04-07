pipeline {
  agent any

  environment {
    DEPLOY_BASE_DIR = '/srv/jianli'
    DEPLOY_ROOT = '/srv/jianli/app-src'
    DEPLOY_ENV_FILE = '/srv/jianli/deploy.env'
    HOST_JENKINS_HOME = '/srv/jenkins_home'
  }

  options {
    timestamps()
  }

  stages {
    stage('Checkout') {
      steps {
        checkout scm
      }
    }

    stage('Sync Deploy Source') {
      steps {
        sh '''
          mkdir -p "$DEPLOY_ROOT"
          rsync -a --delete \
            --exclude '.git' \
            --exclude '.worktrees' \
            --exclude 'web/node_modules' \
            "$WORKSPACE"/ "$DEPLOY_ROOT"/
        '''
      }
    }

    stage('Test Backend') {
      steps {
        sh '''
          HOST_WORKSPACE="${WORKSPACE#/var/jenkins_home}"
          HOST_WORKSPACE="$HOST_JENKINS_HOME$HOST_WORKSPACE"
          docker run --rm \
            -v "$HOST_WORKSPACE":/workspace \
            -w /workspace \
            golang:1.25-bookworm \
            bash -lc "go test ./... -count=1"
        '''
      }
    }

    stage('Test Frontend') {
      steps {
        sh '''
          HOST_WORKSPACE="${WORKSPACE#/var/jenkins_home}"
          HOST_WORKSPACE="$HOST_JENKINS_HOME$HOST_WORKSPACE"
          docker run --rm \
            -v "$HOST_WORKSPACE":/workspace \
            -w /workspace/web \
            node:24-bookworm \
            bash -lc "npm ci && npm run test"
        '''
      }
    }

    stage('Build Images') {
      steps {
        sh '''
          docker compose --env-file "$DEPLOY_ENV_FILE" -f "$DEPLOY_ROOT/deploy/docker-compose.yml" build app nginx
        '''
      }
    }

    stage('Approve Deploy') {
      steps {
        input message: 'Tests and image build passed. Deploy to production?', ok: 'Deploy'
      }
    }

    stage('Deploy') {
      steps {
        sh '''
          chmod +x "$DEPLOY_ROOT/deploy/scripts/deploy.sh"
          "$DEPLOY_ROOT/deploy/scripts/deploy.sh" "$DEPLOY_ROOT" "$DEPLOY_ENV_FILE"
        '''
      }
    }

    stage('Health Check') {
      steps {
        sh '''
          chmod +x "$DEPLOY_ROOT/deploy/scripts/healthcheck.sh"
          APP_HEALTHCHECK_URL=http://127.0.0.1/api/resume "$DEPLOY_ROOT/deploy/scripts/healthcheck.sh"
        '''
      }
    }
  }
}
