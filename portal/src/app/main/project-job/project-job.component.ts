import {Component, Input, OnInit} from '@angular/core';
import {ApiService} from '../../api.service';
import {NzTableQueryParams} from 'ng-zorro-antd';
import {ActivatedRoute, Router} from '@angular/router';
import {TabRef} from '../tabs/tabs.component';

@Component({
  selector: 'app-project-job',
  templateUrl: './project-job.component.html',
  styleUrls: ['./project-job.component.scss']
})
export class ProjectJobComponent implements OnInit {
  @Input() project: any = {};

  jobs: [];
  total = 0;
  pageIndex = 1;
  pageSize = 10;
  sortField = null;
  sortOrder = null;
  filters = [];
  keyword = '';
  loading = false;

  constructor(private as: ApiService, private router: Router, private routeInfo: ActivatedRoute, private tab: TabRef) {
    tab.name = '定时任务';
  }

  ngOnInit(): void {
    this.tab.name = '项目【' + this.project.name + '】定时任务';
  }

  reload(): void {
    this.pageIndex = 1;
    this.keyword = '';
    this.load();
  }

  load(): void {
    this.loading = true;
    this.as.post('/project/' + this.project.id + '/jobs', {
      offset: (this.pageIndex - 1) * this.pageSize,
      length: this.pageSize,
      sortKey: this.sortField,
      sortOrder: this.sortOrder,
      filters: this.filters,
      keyword: this.keyword,
    }).subscribe(res => {

      this.jobs = res.data;
      this.total = res.total;
    }, error => {
      console.log('error', error);
    }, () => {
      this.loading = false;
    });
  }

  create(): void {
    this.router.navigate(['/admin/project/' + this.project.id + '/job/create']);
  }

  edit(c): void {
    this.router.navigate(['/admin/project/' + this.project.id + '/job/' + c.id + '/edit']);
  }

  onTableQuery(params: NzTableQueryParams): void {
    const {pageSize, pageIndex, sort, filter} = params;
    this.pageSize = pageSize;
    this.pageIndex = pageIndex;
    const currentSort = sort.find(item => item.value !== null);
    this.sortField = (currentSort && currentSort.key) || null;
    this.sortOrder = (currentSort && currentSort.value) || null;
    this.filters = filter;
    this.load();
  }

  search(): void {
    this.pageIndex = 1;
    this.load();
  }
}
