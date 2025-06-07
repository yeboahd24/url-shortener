# Deployment Guide 🚀

This guide covers different deployment scenarios for the URL Shortener API.

## 📋 Prerequisites

- Docker and Docker Compose installed
- Domain name (for production)
- SSL certificates (for HTTPS)
- Minimum 1GB RAM, 1 CPU core

## 🏠 Local Development

### Quick Start
```bash
# Clone repository
git clone https://github.com/yeboahd24/url-shortener.git
cd url-shortener

# Setup environment
make dev-setup
# Edit .env file with your settings

# Start development environment
make dev-start

# Run application
make run
```

### Manual Setup
```bash
# Start only database services
docker-compose up postgres redis -d

# Run application locally
go run .
```

## 🌐 Production Deployment

### 1. Server Preparation

**Update System**
```bash
sudo apt update && sudo apt upgrade -y
```

**Install Docker**
```bash
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER
```

**Install Docker Compose**
```bash
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```

### 2. Application Deployment

**Clone and Configure**
```bash
git clone https://github.com/yeboahd24/url-shortener.git
cd url-shortener

# Create production environment file
cp .env.example .env
nano .env
```

**Environment Configuration**
```bash
# Required production settings
POSTGRES_PASSWORD=your_very_secure_password_here
REDIS_PASSWORD=your_redis_password_here
POSTGRES_DSN=postgres://postgres:your_very_secure_password_here@postgres:5432/urlshortener?sslmode=disable
```

**Deploy**
```bash
# Production deployment
make prod-deploy

# Or manually
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```

### 3. SSL Configuration (HTTPS)

**Obtain SSL Certificate**
```bash
# Using Let's Encrypt (recommended)
sudo apt install certbot
sudo certbot certonly --standalone -d your-domain.com
```

**Configure SSL**
```bash
# Create SSL directory
mkdir ssl

# Copy certificates
sudo cp /etc/letsencrypt/live/your-domain.com/fullchain.pem ssl/cert.pem
sudo cp /etc/letsencrypt/live/your-domain.com/privkey.pem ssl/key.pem
sudo chown $USER:$USER ssl/*
```

**Update Nginx Configuration**
```bash
# Edit nginx.conf
nano nginx.conf

# Uncomment HTTPS server block
# Update server_name to your domain
# Restart services
docker-compose restart nginx
```

## ☁️ Cloud Deployment

### AWS EC2

**1. Launch EC2 Instance**
- Instance type: t3.medium or larger
- Security groups: Allow ports 22, 80, 443
- Storage: 20GB+ EBS volume

**2. Setup Application**
```bash
# Connect to instance
ssh -i your-key.pem ubuntu@your-ec2-ip

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Deploy application
git clone https://github.com/yeboahd24/url-shortener.git
cd url-shortener
cp .env.example .env
# Configure .env
make prod-deploy
```

**3. Configure Domain**
- Point your domain to EC2 public IP
- Setup Route 53 (optional)
- Configure SSL with Let's Encrypt

### Google Cloud Platform

**1. Create Compute Engine Instance**
```bash
gcloud compute instances create url-shortener \
  --image-family=ubuntu-2004-lts \
  --image-project=ubuntu-os-cloud \
  --machine-type=e2-medium \
  --tags=http-server,https-server
```

**2. Configure Firewall**
```bash
gcloud compute firewall-rules create allow-http-https \
  --allow tcp:80,tcp:443 \
  --source-ranges 0.0.0.0/0 \
  --target-tags http-server,https-server
```

### DigitalOcean Droplet

**1. Create Droplet**
- Ubuntu 20.04 LTS
- 2GB RAM minimum
- Enable monitoring

**2. Deploy Application**
```bash
# SSH to droplet
ssh root@your-droplet-ip

# Setup and deploy
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh

git clone https://github.com/yeboahd24/url-shortener.git
cd url-shortener
make prod-deploy
```

## 🔧 Configuration Management

### Environment Variables

**Production .env Example**
```bash
# Database
POSTGRES_PASSWORD=super_secure_password_123!
POSTGRES_DSN=postgres://postgres:super_secure_password_123!@postgres:5432/urlshortener?sslmode=disable

# Redis
REDIS_ADDR=redis:6379
REDIS_PASS=redis_secure_password_456!

# Application
PORT=8080
GEO_API_URL=http://ip-api.com/json
API_KEY_HEADER=X-API-Key
GIN_MODE=release

# Domain (optional)
DOMAIN=your-domain.com
```

### Scaling Configuration

**Horizontal Scaling**
```bash
# Scale application instances
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d --scale app=3

# Scale with load balancer
# Update nginx.conf upstream block with multiple servers
```

**Vertical Scaling**
```bash
# Update docker-compose.prod.yml
services:
  app:
    deploy:
      resources:
        limits:
          memory: 512M
        reservations:
          memory: 256M
```

## 📊 Monitoring & Maintenance

### Health Monitoring
```bash
# Check service health
curl http://your-domain.com/health

# Monitor logs
docker-compose logs -f app

# Check resource usage
docker stats
```

### Backup Strategy
```bash
# Database backup
docker-compose exec postgres pg_dump -U postgres urlshortener > backup.sql

# Restore database
docker-compose exec -T postgres psql -U postgres urlshortener < backup.sql
```

### Updates
```bash
# Update application
git pull origin main
docker-compose build app
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```

## 🔒 Security Checklist

- [ ] Strong passwords for database and Redis
- [ ] SSL/TLS certificates configured
- [ ] Firewall rules configured
- [ ] Regular security updates
- [ ] API rate limiting enabled
- [ ] Non-root container execution
- [ ] Environment variables secured
- [ ] Regular backups scheduled

## 🆘 Troubleshooting

### Common Issues

**Service Won't Start**
```bash
# Check logs
docker-compose logs service-name

# Check configuration
docker-compose config

# Restart services
docker-compose restart
```

**Database Connection Issues**
```bash
# Check PostgreSQL status
docker-compose exec postgres pg_isready

# Reset database
docker-compose down -v
docker-compose up -d
```

**SSL Certificate Issues**
```bash
# Renew Let's Encrypt certificate
sudo certbot renew

# Update certificate in container
docker-compose restart nginx
```

## 📞 Support

For deployment issues:
1. Check logs: `docker-compose logs`
2. Verify configuration: `docker-compose config`
3. Test connectivity: `curl http://localhost/health`
4. Create GitHub issue with logs and configuration
