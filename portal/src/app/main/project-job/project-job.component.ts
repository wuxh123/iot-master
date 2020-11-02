import {Component, Input, OnInit} from '@angular/core';
import {ApiService} from '../../api.service';
import {NzModalService, NzTableQueryParams} from 'ng-zorro-antd';
import {Router} from '@angular/router';
import {ProjectJobEditComponent} from '../project-job-edit/project-job-edit.component';

@Component({
  selector: 'app-project-job',
  templateUrl: './project-job.component.html',
  styleUrls: ['./project-job.component.scss']
})
export class ProjectJobComponent implements OnInit {
  @Input() project: any = {};

  inited = false;
  tableQuery: any;

  datum: any[];
  total = 0;
  pageIndex = 1;
  pageSize = 10;
  sortField = null;
  sortOrder = null;
  filters = [];
  keyword = '';
  loading = false;

  constructor(private as: ApiService, private router: Router, private ms: NzModalService) {
  }

  ngOnInit(): void {
    this.inited = true;
    if (this.tableQuery) {
      this.onTableQuery(this.tableQuery);
    }
  }

  reload(): void {
    this.pageIndex = 1;
    this.keyword = '';
    this.load();
  }

  load(): void {
    this.loading = true;
    this.as.post('project/' + this.project.id + '/jobs', {
      offset: (this.pageIndex - 1) * this.pageSize,
      length: this.pageSize,
      sortKey: this.sortField,
      sortOrder: this.sortOrder,
      filters: this.filters,
      keywords: [
        {key: 'Name', value: this.keyword},
      ]
    }).subscribe(res => {

      this.datum = res.data;
      this.total = res.total;
    }, error => {
      console.log('error', error);
    }, () => {
      this.loading = false;
    });
  }


  edit(id?): void {
    const modal = this.ms.create({
      nzTitle: id ? '编辑任务' : '创建任务',
      nzContent: ProjectJobEditComponent,
      nzFooter: null,
      nzMaskClosable: false,
      // nzViewContainerRef: this.viewContainerRef,
      nzComponentParams: {id},
    });
    // insert/update after close
    modal.afterClose.subscribe(data => {
      if (!data) {
        return;
      }
      if (id) {
        this.datum.forEach((c: any, i, a: any[]) => {
            if (c.id === data.id) {
              a[i] = data;
            }
          }
        );
      } else {
        this.datum.unshift(data);
      }
    });
  }


  onTableQuery(params: NzTableQueryParams): void {
    if (!this.inited) {
      this.tableQuery = params;
      return;
    }
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
