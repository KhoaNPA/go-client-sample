
Part 1: Architectural Challenge
- Q1:How would you improve the current design to achieve better: 
    High availability, Resilience, Performance, Cost efficiency
    
Answer: The anwser is in general because Workload/Usage isn't mentioned
- At step (1)(2), insteads of only pushing the ID of Scan entity into rabbitmq. API(1) calls StorageAPI to create a Scan entity in the DB and gets related data also. => We can remove (4) to save time calling StorageAPI. As our API Services and EKS Cluster stay in the same VPC and apply VPC Endpoints, larger payload size for each RabbitMQ message isn't a problem.
- At (5), insteads of create Kubernetes Job with sourcecode, use AWS NFS to store Git Repository locally and mount NFS to Kubernetes jobs. So new launched jobs can checkout repository data faster.
- Each job takes 5s-30mins, so we can consider using EKS Node group with EC2 Spot instances - Save cost, and node spinup time is not a problem. Use multiple GitHub Service accounts to prevent Github rate limit.
- Monitor EC2 usage to consider buying EC2 Reserved Instance (EC2 RI) to save cost.
- Docker images can be stored in AWS ECR and we apply ECR endpoint to reduce latency when EKS pulls images.
- High availability for RDS: RDS Read-replica, Multi-AZ, Backup RDS snapshot frequently.
- Use Redis or Memcache to reduce latency, improve read performance and reduce stress for RDS
- Consider putting StorageAPI (or all API services) to EKS too to have High Avalability & Resilience, because StorageAPI is handling many requests from many other sources.
- Store static data of Dashboard in AWS Cloudfront (or any CDN)

- Q2: The number of scan requests can increase/decrease randomly in a day and on most weekends the system receives almost no requests at all.
    - Use EC2 Reserved Instance (EC2 RI)
    - Apply EKS Node group with EC2 Spot instance
    - Monitor Scan Usage to define EC2 Auto Scaling Group (ASG) Min-Desired-Max Capacity with suitable value. As we know weekends usage, create a Lambda function (to make sure ASG is always updated) to adjust ASG with different values for Weekdays and Weekends.
    - As we use EKS cluster, we may set ASG Min equals 0 at weekends, when new jobs are scheduled, EKS can request ASG to launch new ec2s but may increase job complete period.

- Q3: In step 6, each job needs to mount the source code folder into every engine that needs to run. How would you store the source code and make sure that engines can run in a scalable way?
    - Consider using AWS NFS to store repositories data locally and shares across all jobs, when job can checkout repository faster as it only need to download new data insteads of entire reposioty. 
    - Use Git Shallow Clone to checkout only 1 branch to download data faster.

- Q4: Propose a high-level disaster recovery (DR) plan for the current architecture.
    Architecture design:
    - ALWAYS put VPC, EKS, ELB, ASG, RDS Subnet group, Lambda... to ALL AZs. Use RabbitMQ Cluster for High Availability and spreads to all AZs
    - For RDS, Use RDS Multi-AZ, Read-replica RDS and frequently take snapshot to protect data.
    - As we use EKS, NFS, EC2 ASG, ELB, so we don't need to worry about EKS clusters and its data when any AZ is down.
    - Use Lambda function to manage ASG.
    - Use ECR to store docker images.
    - Store static data in dashboard to S3, Cloudfront
    - Consider put all API services to EKS or Lambda if possible

    DR Plan:
    - Practice DR 1-2 times every year in weekend.
    - Select 1 random AZ to be removed for simulating AZ down.
    - Monitor system metrics from beginning.
    - DR Simulation:
        - Create system metric reports before running DR simulation
        - Temporarily turn-off all alerts. (ignore if you want to test alerting system)
        - Remove the above AZ from ASG, ELB, Lambda subnet, RDS subnet group and RabbitMQ Node. Ensure the Service API, Lamnda, scan job are scheduled to other AZs and RabbitMQ is running.
        - Force-failover for RDS and ensure RDS is started again in other AZs and has correct data. StorageAPI must have re-create DB connection function when current connection is dropped, or refer to Redis/Memcache/Read-replica when Master RDS is unavailable.
        - Run system test to verify everything is up and running, data is assurance.
    - Revert DR simulation:
        - Add the removed AZ back to ASG, ELB, Lambda subnet, RDS subnet group and RabbitMQ Cluster.
        - Execute force-failover for RDS again (can be skipped).
        - Turn-on all alerts back.
        - Run system test to verify everything is up and running, data is assurance.
        - Create system metric reports after running DR simulation
    - Check monitor system.
    - Compare Pre-DR reports vs Post-DR reports.

---
Part 2: Technical Challenge
